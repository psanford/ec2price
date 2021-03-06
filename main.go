package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

var (
	basePriceURL = "https://pricing.us-east-1.amazonaws.com"
	indexPath    = "/offers/v1.0/aws/index.json"

	region           = flag.String("region", "us-east-1", "AWS Region")
	fetchOffers      = flag.Bool("fetch_offers", false, "Fetch offers and price file to disk")
	familyTypes      = flag.Bool("family", false, "Print family type information")
	checkFamilyTypes = flag.Bool("check_family", false, "Check family types against instance types (for missing types)")
	csvOutput        = flag.Bool("csv", false, "output as csv")
)

func main() {
	flag.Parse()

	if *familyTypes {
		seen := make(map[string]bool)
		for _, ft := range instanceTypes {
			if seen[ft.Name] {
				fmt.Fprintf(os.Stderr, "!!! Duplicate family type found: %s\n", ft)
			}
			seen[ft.Name] = true
			fmt.Printf("%5.5s %4d %10.10s %s\n", ft.Name, ft.Year, ft.Prefix, ft.Flags)
		}
		return
	}

	r, err := http.Get(basePriceURL + indexPath)
	checkErr(err, "Get index")

	if r.StatusCode != 200 {
		b, _ := ioutil.ReadAll(r.Body)
		log.Fatalf("Get index status %d\n%s", r.StatusCode, b)
	}

	br := teeToFile(r.Body, "/tmp/ec2-price-index.json")
	dec := json.NewDecoder(br)
	var idx PriceIndex
	err = dec.Decode(&idx)
	checkErr(err, "Read index json")
	br.Close()

	regionURL := idx.Offers["AmazonEC2"].CurrentRegionIndexURL

	r, err = http.Get(basePriceURL + regionURL)
	checkErr(err, "Get region index")

	if r.StatusCode != 200 {
		b, _ := ioutil.ReadAll(r.Body)
		log.Fatalf("Get region index status %d\n%s", r.StatusCode, b)
	}

	br = teeToFile(r.Body, "/tmp/ec2-price-region-index.json")
	dec = json.NewDecoder(br)
	var regionIdx RegionIndex
	err = dec.Decode(&regionIdx)
	checkErr(err, "Read region index json")
	br.Close()

	ec2Path := regionIdx.Regions[*region].CurrentVersionURL

	r, err = http.Get(basePriceURL + ec2Path)
	checkErr(err, "Get prices")

	if r.StatusCode != 200 {
		b, _ := ioutil.ReadAll(r.Body)
		log.Fatalf("Get prices status %d\n%s", r.StatusCode, b)
	}

	br = teeToFile(r.Body, "/tmp/ec2-price.json")
	dec = json.NewDecoder(br)
	var prices PriceDoc
	err = dec.Decode(&prices)
	checkErr(err, "Read price json")
	br.Close()

	var instances []InstanceType
	families := make(map[string]familyInfo)

	for sku, prod := range prices.Products {
		attrs := prod.Attributes

		if strings.Index(attrs.InstanceType, ".") == -1 ||
			!strings.HasPrefix(attrs.UsageType, "BoxUsage:") ||
			attrs.OperatingSystem != "Linux" ||
			attrs.Operation != "RunInstances" {
			continue
		}

		var reservedAnnual float64
		skuTerms := prices.Terms.Reserved[sku]
	RESERVATION:
		for _, term := range skuTerms {
			if term.TermAttributes.LeaseContractLength == "1yr" &&
				term.TermAttributes.PurchaseOption == "No Upfront" &&
				term.TermAttributes.OfferingClass == "convertible" {
				for _, pd := range term.PriceDimensions {
					if pd.Unit == "Hrs" {
						f, _ := strconv.ParseFloat(pd.PricePerUnit["USD"], 64)
						reservedAnnual = f * 24.0 * 365.0
						break RESERVATION
					}
				}
				break RESERVATION
			}
		}

		var onDemandCost float64
		var hourly float64
		onDemand := prices.Terms.OnDemand[sku]
	ONDEMANDOUTER:
		for _, od := range onDemand {
			for _, pd := range od.PriceDimensions {
				f, _ := strconv.ParseFloat(pd.PricePerUnit["USD"], 64)
				hourly = f
				onDemandCost = f * 24.0 * 365.0
				break ONDEMANDOUTER
			}
		}

		memS := strings.TrimSuffix(attrs.Memory, " GiB")
		memS = strings.ReplaceAll(memS, ",", "")
		mem, _ := strconv.ParseFloat(memS, 64)

		diskTotal := 0
		if attrs.Storage != "EBS only" {
			parts := strings.Split(attrs.Storage, " ")
			if len(parts) == 1 {
				diskTotal, _ = strconv.Atoi(parts[0])
			} else if len(parts) > 2 {
				if parts[1] == "x" {
					count, _ := strconv.Atoi(parts[0])
					size, _ := strconv.Atoi(strings.ReplaceAll(parts[2], ",", ""))
					diskTotal = count * size
				}

			}
		}

		instance := InstanceType{
			Name:           attrs.InstanceType,
			VCPU:           attrs.VCPU,
			Memory:         mem,
			Disk:           strings.Replace(attrs.Storage, " SSD", "", 1),
			DiskTotal:      diskTotal,
			Hourly:         hourly,
			OnDemandAnnual: onDemandCost,
			ReservedAnnual: reservedAnnual,
			CPUMfgr:        mfgrFromString(attrs.PhysicalProcessor),
		}

		family := strings.SplitN(attrs.InstanceType, ".", 2)[0]
		families[family] = familyInfo{
			InstanceFamily:    attrs.InstanceFamily,
			PhysicalProcessor: attrs.PhysicalProcessor,
			CPUMfgr:           instance.CPUMfgr,
			CurrentGen:        attrs.CurrentGeneration == "Yes",
		}

		instances = append(instances, instance)
	}

	sort.Slice(instances, func(a, b int) bool { return instances[a].OnDemandAnnual < instances[b].OnDemandAnnual })

	fieldNames := []string{"type", "mem", "vcpu", "disk", "dsk", "mfg", "hourly", "annual", "annual-reserved"}

	if *csvOutput {
		w := csv.NewWriter(os.Stdout)
		w.Write(fieldNames)
		for _, in := range instances {
			w.Write([]string{
				toS(in.Name),
				toS(in.Memory),
				toS(in.VCPU),
				toS(in.Disk),
				toS(in.DiskTotal),
				toS(in.CPUMfgr),
				toS(in.Hourly),
				toS(in.OnDemandAnnual),
				toS(in.ReservedAnnual),
			})
		}
		w.Flush()
		return
	}

	format := "%15s %10.01f %6s %15s %5d %3s %9.04f %9.02f %.2f\n"
	var fieldNamesI []interface{} = make([]interface{}, len(fieldNames))
	for i, d := range fieldNames {
		fieldNamesI[i] = d
	}
	fmt.Printf("%15s %10s %6s %15s %5s %3s %9s %9s %s\n", fieldNamesI...)

	for _, in := range instances {
		fmt.Printf(format, in.Name, in.Memory, in.VCPU, in.Disk, in.DiskTotal, in.CPUMfgr, in.Hourly, in.OnDemandAnnual, in.ReservedAnnual)
	}

	if *checkFamilyTypes {

		familyNames := make([]string, 0, len(families))

		for family := range families {
			familyNames = append(familyNames, family)
		}

		sort.Strings(familyNames)

		for _, family := range familyNames {
			var found bool
			for _, it := range instanceTypes {
				if it.Name == family {
					found = true
					fmt.Println(it.String())
					break
				}
			}
			if !found {
				fmt.Printf("!!!missing %s\n", family)
			}
		}
	}
}

type InstanceTypeInfo struct {
	Name   string
	Year   int
	Prefix InstanceCodePrefix
	Flags  InstanceCodeSuffix
}

func (it InstanceTypeInfo) String() string {
	return fmt.Sprintf("%s %d %s %s", it.Name, it.Year, it.Prefix, it.Flags)
}

var instanceTypes = []InstanceTypeInfo{
	{
		Name:   "m1",
		Year:   2006,
		Prefix: MainPrefix,
	},
	{
		Name:   "c1",
		Year:   2008,
		Prefix: CpuPrefix,
	},
	{
		Name:   "m2",
		Year:   2009,
		Prefix: MainPrefix,
	},
	{
		Name:   "cc1",
		Year:   2010,
		Prefix: ClusterComputePrefix,
	},
	{
		Name:   "t1",
		Year:   2010,
		Prefix: BurstPrefix,
	},
	{
		Name:   "cg1",
		Year:   2010,
		Prefix: GPUPrefix,
	},
	{
		Name:   "cc2",
		Year:   2011,
		Prefix: ClusterComputePrefix,
	},
	{
		Name:   "hi1",
		Year:   2012,
		Prefix: SSDPrefix,
	},
	{
		Name:   "m3",
		Year:   2012,
		Prefix: MainPrefix,
	},
	{
		Name:   "hs1",
		Year:   2012,
		Prefix: DenseHDDPrefix,
	},
	{
		Name:   "cr1",
		Year:   2013,
		Prefix: ClusterComputePrefix,
	},
	{
		Name:   "c3",
		Year:   2013,
		Prefix: CpuPrefix,
	},
	{
		Name:   "g2",
		Year:   2013,
		Prefix: GPUPrefix,
	},
	{
		Name:   "i2",
		Year:   2013,
		Prefix: SSDPrefix,
	},
	{
		Name:   "r3",
		Year:   2014,
		Prefix: MemMorePrefix,
	},
	{
		Name:   "t2",
		Year:   2014,
		Prefix: BurstPrefix,
	},
	{
		Name:   "c4",
		Year:   2015,
		Prefix: CpuPrefix,
	},
	{
		Name:   "d2",
		Year:   2015,
		Prefix: DenseHDDPrefix,
	},
	{
		Name:   "m4",
		Year:   2015,
		Prefix: MainPrefix,
	},
	{
		Name:   "x1",
		Year:   2016,
		Prefix: MemXtremePrefix,
	},
	{
		Name:   "p2",
		Year:   2016,
		Prefix: GPUPrefix,
	},
	{
		Name:   "f1",
		Year:   2016,
		Prefix: FPGAPrefix,
	},
	{
		Name:   "r4",
		Year:   2016,
		Prefix: MemMorePrefix,
	},
	{
		Name:   "i3",
		Year:   2016,
		Prefix: SSDPrefix,
	},
	{
		Name:   "c5",
		Year:   2016,
		Prefix: CpuPrefix,
	},
	{
		Name:   "g3",
		Year:   2017,
		Prefix: GPUPrefix,
	},
	{
		Name:   "x1e",
		Year:   2017,
		Prefix: MemXtremePrefix,
	},
	{
		Name:   "p3",
		Year:   2017,
		Prefix: GPUPrefix,
	},
	{
		Name:   "m5",
		Year:   2017,
		Prefix: MainPrefix,
	},
	{
		Name:   "h1",
		Year:   2017,
		Prefix: DenseHDDPrefix,
	},
	{
		Name:   "c5d",
		Year:   2018,
		Prefix: CpuPrefix,
		Flags:  NVMeSuffix,
	},
	{
		Name:   "m5d",
		Year:   2018,
		Prefix: MainPrefix,
		Flags:  NVMeSuffix,
	},
	{
		Name:   "z1d",
		Year:   2018,
		Prefix: HighFreqPrefix,
	},
	{
		Name:   "r5",
		Year:   2018,
		Prefix: MemMorePrefix,
	},
	{
		Name:   "t3",
		Year:   2018,
		Prefix: BurstPrefix,
	},
	{
		Name:   "g3s",
		Year:   2018,
		Prefix: GPUPrefix,
	},
	{
		Name:   "m5a",
		Year:   2018,
		Prefix: MainPrefix,
		Flags:  AmdSuffix,
	},
	{
		Name:   "r5a",
		Year:   2018,
		Prefix: MemMorePrefix,
		Flags:  AmdSuffix,
	},
	{
		Name:   "c5n",
		Year:   2018,
		Prefix: CpuPrefix,
		Flags:  NetworkSuffix,
	},
	{
		Name:   "a1",
		Year:   2018,
		Prefix: ArmPrefix,
	},
	{
		Name:   "p3dn",
		Year:   2018,
		Prefix: GPUPrefix,
		Flags:  NVMeSuffix | NetworkSuffix,
	},
	{
		Name:   "g4",
		Year:   2019,
		Prefix: GPUPrefix,
	},
	{
		Name:   "m5ad",
		Year:   2019,
		Prefix: MainPrefix,
		Flags:  AmdSuffix | NVMeSuffix,
	},
	{
		Name:   "r5d",
		Year:   2019,
		Prefix: MemMorePrefix,
		Flags:  NVMeSuffix,
	},
	{
		Name:   "r5ad",
		Year:   2019,
		Prefix: MemMorePrefix,
		Flags:  AmdSuffix | NVMeSuffix,
	},
	{
		Name:   "i3en",
		Year:   2019,
		Prefix: SSDPrefix,
		Flags:  NetworkSuffix,
	},
	{
		Name:   "g4dn",
		Year:   2019,
		Prefix: GPUPrefix,
		Flags:  GpuAmdSuffix,
	},
	{
		Name:   "r5dn",
		Year:   2019,
		Prefix: MemMorePrefix,
		Flags:  NVMeSuffix | NetworkSuffix,
	},
	{
		Name:   "r5n",
		Year:   2019,
		Prefix: MemMorePrefix,
		Flags:  NetworkSuffix,
	},
	{
		Name:   "m5dn",
		Year:   2019,
		Prefix: MainPrefix,
		Flags:  NVMeSuffix | NetworkSuffix,
	},
	{
		Name:   "m5n",
		Year:   2019,
		Prefix: MainPrefix,
		Flags:  NetworkSuffix,
	},
	{
		Name:   "inf1",
		Year:   2019,
		Prefix: InferencePrefix,
	},
	{
		Name:   "t3a",
		Year:   2019,
		Prefix: BurstPrefix,
		Flags:  AmdSuffix,
	},
	{
		Name:   "c5a",
		Year:   2020,
		Prefix: CpuPrefix,
		Flags:  AmdSuffix,
	},
	{
		Name:   "c5ad",
		Year:   2020,
		Prefix: CpuPrefix,
		Flags:  AmdSuffix | NVMeSuffix,
	},
	{
		Name:   "c6g",
		Year:   2020,
		Prefix: CpuPrefix,
		Flags:  GravitonSuffix,
	},
	{
		Name:   "c6gn",
		Year:   2020,
		Prefix: CpuPrefix,
		Flags:  GravitonSuffix | NetworkSuffix,
	},
	{
		Name:   "c6gd",
		Year:   2020,
		Prefix: CpuPrefix,
		Flags:  GravitonSuffix | NVMeSuffix,
	},
	{
		Name:   "d3",
		Year:   2020,
		Prefix: DenseHDDPrefix,
	},
	{
		Name:   "d3en",
		Year:   2020,
		Prefix: DenseHDDPrefix,
		Flags:  NetworkSuffix,
	},
	{
		Name:   "g4ad",
		Year:   2020,
		Prefix: GPUPrefix,
		Flags:  AmdSuffix | NVMeSuffix,
	},
	{
		Name:   "m5zn",
		Year:   2020,
		Prefix: MainPrefix,
		Flags:  HighFreqSuffix | NetworkSuffix,
	},
	{
		Name:   "m6g",
		Year:   2020,
		Prefix: MainPrefix,
		Flags:  GravitonSuffix,
	},
	{
		Name:   "m6gn",
		Year:   2020,
		Prefix: MainPrefix,
		Flags:  GravitonSuffix | NetworkSuffix,
	},
	{
		Name:   "p4d",
		Year:   2020,
		Prefix: GPUPrefix,
		Flags:  NVMeSuffix,
	},
	{
		Name:   "r5b",
		Year:   2020,
		Prefix: MemMorePrefix,
		Flags:  EBSOptimizedSuffix,
	},
	{
		Name:   "r6g",
		Year:   2020,
		Prefix: MemMorePrefix,
		Flags:  GravitonSuffix,
	},
	{
		Name:   "r6gd",
		Year:   2020,
		Prefix: MemMorePrefix,
		Flags:  GravitonSuffix | NVMeSuffix,
	},
	{
		Name:   "t4g",
		Year:   2020,
		Prefix: BurstPrefix,
		Flags:  GravitonSuffix,
	},
}

type InstanceCodePrefix int

const (
	ArmPrefix InstanceCodePrefix = iota
	BurstPrefix
	MainPrefix
	CpuPrefix
	MemMorePrefix
	MemXtremePrefix
	HighFreqPrefix
	GPUPrefix
	InferencePrefix
	FPGAPrefix
	SSDPrefix
	DenseHDDPrefix
	ClusterComputePrefix
)

type InstanceCodeSuffix int

const (
	GravitonSuffix InstanceCodeSuffix = 1 << iota
	AmdSuffix
	NVMeSuffix
	NetworkSuffix
	GpuNvidiaSuffix
	GpuAmdSuffix
	HighFreqSuffix
	EBSOptimizedSuffix
)

func (c InstanceCodeSuffix) String() string {
	var parts []string

	if c&GravitonSuffix == GravitonSuffix {
		parts = append(parts, "graviton")
	}
	if c&AmdSuffix == AmdSuffix {
		parts = append(parts, "amd")
	}
	if c&NVMeSuffix == NVMeSuffix {
		parts = append(parts, "nvme")
	}
	if c&NetworkSuffix == NetworkSuffix {
		parts = append(parts, "net")
	}
	if c&GpuNvidiaSuffix == GpuNvidiaSuffix {
		parts = append(parts, "gpu-nvidia")
	}
	if c&GpuAmdSuffix == GpuAmdSuffix {
		parts = append(parts, "gpu-amd")
	}
	if c&HighFreqSuffix == HighFreqSuffix {
		parts = append(parts, "high-freq")
	}
	if c&EBSOptimizedSuffix == EBSOptimizedSuffix {
		parts = append(parts, "ebs-optimized")
	}

	return strings.Join(parts, ",")
}

func (c InstanceCodePrefix) String() string {
	switch c {
	case ArmPrefix:
		return "arm"
	case BurstPrefix:
		return "burst"
	case MainPrefix:
		return "main"
	case CpuPrefix:
		return "cpu"
	case MemMorePrefix:
		return "more-mem"
	case MemXtremePrefix:
		return "mem-xtreme"
	case HighFreqPrefix:
		return "high-freq"
	case GPUPrefix:
		return "gpu"
	case InferencePrefix:
		return "inference"
	case FPGAPrefix:
		return "fpga"
	case SSDPrefix:
		return "ssd"
	case DenseHDDPrefix:
		return "dense-hdd"
	case ClusterComputePrefix:
		return "cluster-compute"
	}
	return fmt.Sprintf("unknown<%x>", int(c))
}

func teeToFile(rc io.ReadCloser, path string) io.ReadCloser {
	if !*fetchOffers {
		return rc
	}

	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}

	tr := teeReader{
		Reader: io.TeeReader(rc, f),
		rc:     rc,
		f:      f,
		name:   path,
	}

	return &tr
}

type teeReader struct {
	io.Reader

	rc   io.ReadCloser
	f    *os.File
	name string
}

func (tr *teeReader) Close() error {
	tr.f.Close()
	fmt.Printf("wrote %s\n", tr.name)
	return tr.rc.Close()
}

type CPUManufacturer int

const (
	CPUIntel CPUManufacturer = 1
	CPUAMD   CPUManufacturer = 2
	CPUAWS   CPUManufacturer = 3
)

func (c CPUManufacturer) String() string {
	switch c {
	case CPUIntel:
		return "int"
	case CPUAMD:
		return "amd"
	case CPUAWS:
		return "arm"
	}

	return "unk"
}

func mfgrFromString(s string) CPUManufacturer {
	if strings.Contains(s, "Intel") {
		return CPUIntel
	} else if strings.Contains(s, "AMD") {
		return CPUAMD
	} else if strings.Contains(s, "AWS") {
		return CPUAWS
	}
	return 0
}

type InstanceType struct {
	Name           string
	VCPU           string
	Memory         float64
	SSD            bool
	Disk           string
	DiskTotal      int
	Hourly         float64
	OnDemandAnnual float64
	ReservedAnnual float64
	CPUMfgr        CPUManufacturer
	CurrentGen     string
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalf("Error: %s: %s", msg, err)
	}
}

type PriceIndex struct {
	Disclaimer      string           `json:"disclaimer"`
	FormatVersion   string           `json:"formatVersion"`
	Offers          map[string]Offer `json:"offers"`
	PublicationDate string           `json:"publicationDate"`
}

type Offer struct {
	CurrentRegionIndexURL string `json:"currentRegionIndexUrl"`
	CurrentVersionURL     string `json:"currentVersionUrl"`
	OfferCode             string `json:"offerCode"`
	VersionIndexURL       string `json:"versionIndexUrl"`
}

type RegionIndex struct {
	Disclaimer      string `json:"disclaimer"`
	FormatVersion   string `json:"formatVersion"`
	PublicationDate string `json:"publicationDate"`
	Regions         map[string]struct {
		CurrentVersionURL string `json:"currentVersionUrl"`
		RegionCode        string `json:"regionCode"`
	} `json:"regions"`
}

type PriceDoc struct {
	Disclaimer      string             `json:"disclaimer"`
	FormatVersion   string             `json:"formatVersion"`
	PublicationDate string             `json:"publicationDate"`
	OfferCode       string             `json:"offerCode"`
	Products        map[string]Product `json:"products"`
	Terms           struct {
		OnDemand map[string]map[string]Term `json:"OnDemand"`
		Reserved map[string]map[string]Term `json:"Reserved"`
	} `json:"terms"`
	Version string `json:"version"`
}

type Term struct {
	EffectiveDate   string `json:"effectiveDate"`
	OfferTermCode   string `json:"offerTermCode"`
	PriceDimensions map[string]struct {
		AppliesTo    []interface{}     `json:"appliesTo"`
		BeginRange   string            `json:"beginRange"`
		Description  string            `json:"description"`
		EndRange     string            `json:"endRange"`
		PricePerUnit map[string]string `json:"pricePerUnit"`
		RateCode     string            `json:"rateCode"`
		Unit         string            `json:"unit"`
	} `json:"priceDimensions"`
	Sku            string `json:"sku"`
	TermAttributes struct {
		LeaseContractLength string `json:"LeaseContractLength"`
		OfferingClass       string `json:"OfferingClass"`
		PurchaseOption      string `json:"PurchaseOption"`
	} `json:"termAttributes"`
}

type Product struct {
	Attributes struct {
		CapacityStatus              string `json:"capacitystatus"`
		ClockSpeed                  string `json:"clockSpeed"`
		CurrentGeneration           string `json:"currentGeneration"`
		DedicatedEBSThroughput      string `json:"dedicatedEbsThroughput"`
		ECU                         string `json:"ecu"`
		EnhancedNetworkingSupported string `json:"enhancedNetworkingSupported"`
		GPU                         string `json:"gpu"`
		InstanceFamily              string `json:"instanceFamily"`
		InstanceType                string `json:"instanceType"`
		IntelAVX2Available          string `json:"intelAvx2Available"`
		IntelAVXAvailable           string `json:"intelAvxAvailable"`
		IntelTurboAvailable         string `json:"intelTurboAvailable"`
		LicenseModel                string `json:"licenseModel"`
		Location                    string `json:"location"`
		LocationType                string `json:"locationType"`
		Memory                      string `json:"memory"`
		NetworkPerformance          string `json:"networkPerformance"`
		NormalizationSizeFactor     string `json:"normalizationSizeFactor"`
		OperatingSystem             string `json:"operatingSystem"`
		Operation                   string `json:"operation"`
		PhysicalProcessor           string `json:"physicalProcessor"`
		PreInstalledSW              string `json:"preInstalledSw"`
		ProcessorArchitecture       string `json:"processorArchitecture"`
		ProcessorFeatures           string `json:"processorFeatures"`
		ServiceCode                 string `json:"servicecode"`
		ServiceName                 string `json:"servicename"`
		Storage                     string `json:"storage"`
		Tenancy                     string `json:"tenancy"`
		UsageType                   string `json:"usagetype"`
		VCPU                        string `json:"vcpu"`
	} `json:"attributes"`
	ProductFamily string `json:"productFamily"`
	Sku           string `json:"sku"`
}

type familyInfo struct {
	InstanceFamily    string
	PhysicalProcessor string
	CPUMfgr           CPUManufacturer
	CurrentGen        bool
}

func toS(i interface{}) string {
	switch v := i.(type) {
	case bool:
		return strconv.FormatBool(v)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', 3, 64)
	case float64:
		return strconv.FormatFloat(v, 'f', 3, 64)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case string:
		return v
	case fmt.Stringer:
		return v.String()
	default:
		return fmt.Sprintf("%+v", i)
	}
}
