package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var (
	basePriceURL = "https://pricing.us-east-1.amazonaws.com"
	indexPath    = "/offers/v1.0/aws/index.json"

	region           = flag.String("region", "us-east-1", "AWS Region")
	fetchOffers      = flag.Bool("fetch-offers", false, "Fetch offers and price file to disk")
	familyTypes      = flag.Bool("family", false, "Print family type information")
	checkFamilyTypes = flag.Bool("check-family", false, "Check family types against instance types (for missing types)")
	outFormat        = flag.String("format", "col", "output format: (col|csv|json")
	shortTypes       = flag.Bool("short-type", false, "output using short type names")
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
		b, _ := io.ReadAll(r.Body)
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
		b, _ := io.ReadAll(r.Body)
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
		b, _ := io.ReadAll(r.Body)
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

		disk, err := parseStorage(attrs.Storage)
		if err != nil {
			log.Printf("parse storage for %s err: %s", attrs.InstanceType, err)
		}

		instType := attrs.InstanceType
		if *shortTypes {
			instType = shortType(instType)
		}

		np, err := parseNetPerf(attrs.NetworkPerformance)
		if err != nil {
			log.Print(err)
		}

		instance := InstanceType{
			Name:           instType,
			VCPU:           attrs.VCPU,
			Memory:         mem,
			Disk:           disk,
			Hourly:         hourly,
			OnDemandAnnual: onDemandCost,
			ReservedAnnual: reservedAnnual,
			CPUMfgr:        mfgrFromString(attrs.PhysicalProcessor),
			NetworkPerf:    np,
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

	fieldNames := []string{"type", "mem", "vcpu", "disk", "mfg", "net", "hourly", "annual", "annual-reserved"}

	if *outFormat == "csv" {
		w := csv.NewWriter(os.Stdout)
		w.Write(fieldNames)
		for _, in := range instances {
			w.Write([]string{
				toS(in.Name),
				toS(in.Memory),
				toS(in.VCPU),
				toS(in.Disk),
				toS(in.CPUMfgr),
				toS(in.NetworkPerf),
				toS(in.Hourly),
				toS(in.OnDemandAnnual),
				toS(in.ReservedAnnual),
			})
		}
		w.Flush()
		return
	} else if *outFormat == "json" {
		w := json.NewEncoder(os.Stdout)
		w.SetIndent("", "  ")
		for _, in := range instances {
			w.Encode(in)
		}
		return
	}
	format := "%17s %10.01f %6s %15s %3s %6s %9.04f %9.02f %.2f\n"
	var fieldNamesI []interface{} = make([]interface{}, len(fieldNames))
	for i, d := range fieldNames {
		fieldNamesI[i] = d
	}
	fmt.Printf("%17s %10s %6s %15s %3s %6s %9s %9s %s\n", fieldNamesI...)

	for _, in := range instances {
		fmt.Printf(format, in.Name, in.Memory, in.VCPU, in.Disk, in.CPUMfgr, in.NetworkPerf, in.Hourly, in.OnDemandAnnual, in.ReservedAnnual)
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

// instanceTypes is generated from families.ndjson into families.go by
// generate_families.go. To add a family, append a line to families.ndjson and
// run `go generate`.
//
//go:generate go run generate_families.go

type InstanceCodePrefix int

const (
	ArmPrefix              InstanceCodePrefix = iota // a
	BurstPrefix                                      // t
	MainPrefix                                       // m
	CpuPrefix                                        // c
	MemMorePrefix                                    // r
	MemXtremePrefix                                  // x
	HighFreqPrefix                                   // z
	GPUPrefix                                        // p,g
	InferencePrefix                                  // inf,trn("Trainium")
	FPGAPrefix                                       // f
	SSDPrefix                                        // i
	DenseHDDPrefix                                   // d
	ClusterComputePrefix                             // cc,cr
	VideoTranscodingPrefix                           // vt
	HPCPrefix                                        // hpc
	XeonScalablePrefix                               // x2
	MemUltraPrefix                                   // u
)

type InstanceCodeSuffix int

const (
	GravitonSuffix       InstanceCodeSuffix = 1 << iota // 'g'
	AmdSuffix                                           // 'a'
	NVMeSuffix                                          // 'd'
	NetworkSuffix                                       // 'n'
	GpuNvidiaSuffix                                     // 'dn'
	GpuAmdSuffix                                        // 'ad'
	HighFreqSuffix                                      // 'z'
	EBSOptimizedSuffix                                  // 'b'
	IntelSuffix                                         // 'i'
	ExtendedMemorySuffix                                // 'e' in memory optimized families (https://aws.amazon.com/blogs/aws/new-amazon-ec2-x2iezn-instances-powered-by-the-fastest-intel-xeon-scalable-cpu-for-memory-intensive-workloads/)
	FlexSuffix                                          // 'flex'
	HpeSuffix                                           // 'h'
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
	if c&IntelSuffix == IntelSuffix {
		parts = append(parts, "intel")
	}
	if c&ExtendedMemorySuffix == ExtendedMemorySuffix {
		parts = append(parts, "extend-mem")
	}
	if c&FlexSuffix == FlexSuffix {
		parts = append(parts, "flex")
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
	case VideoTranscodingPrefix:
		return "video-transcode"
	case HPCPrefix:
		return "hpc"
	case XeonScalablePrefix:
		return "xeon"
	case MemUltraPrefix:
		return "mem-ultra"
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
	Disk           Disk
	Hourly         float64
	OnDemandAnnual float64
	ReservedAnnual float64
	CPUMfgr        CPUManufacturer
	CurrentGen     string
	NetworkPerf    NetworkPerf
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

var typeReplacer = strings.NewReplacer(
	"large", "l",
	"medium", "m",
	"metal", "⛁",
	"micro", "μ",
	"nano", "n",
	"small", "s",
)

func shortType(fullType string) string {
	return typeReplacer.Replace(fullType)
}

type NetworkPerf struct {
	CapGb    float64
	Bursting bool
}

func (np NetworkPerf) String() string {
	var burstIndicator string

	if np.Bursting {
		burstIndicator = "*"
	}
	return fmt.Sprintf("%0.1f%s", np.CapGb, burstIndicator)
}

var netPerfRE = regexp.MustCompile(`(Up to )?(\d+) (Gigabit|Megabit)`)

func parseNetPerf(n string) (NetworkPerf, error) {
	var perf NetworkPerf

	m := netPerfRE.FindStringSubmatch(n)
	if len(m) > 0 {

		nStr := m[2]
		f, _ := strconv.ParseFloat(nStr, 64)
		if m[3] == "Megabit" {
			f = f / 1000
		}

		perf.CapGb = f
		if m[1] != "" {
			perf.Bursting = true
		}

		return perf, nil
	}

	words := map[string]NetworkPerf{
		"Very Low": {
			CapGb:    0.01,
			Bursting: true,
		},
		"High": {
			CapGb: 1,
		},
		"Low": {
			CapGb:    0.01,
			Bursting: true,
		},
		"Low to Moderate": {
			CapGb:    0.01,
			Bursting: true,
		},
		"Moderate": {
			CapGb:    0.1,
			Bursting: true,
		},
		"NA": {
			CapGb: 1,
		},
	}

	if match, found := words[n]; found {
		return match, nil
	}

	return NetworkPerf{}, fmt.Errorf("failed to parse network perf: %q", n)
}

type Disk struct {
	Count     int // count 0 means EBSOnly
	PerDiskGB int
	SSD       bool
	NVMe      bool
}

func (d Disk) String() string {
	if d.Count == 0 {
		return "EBS"
	}
	suffix := "GB"
	total := d.Count * d.PerDiskGB
	if total > 1000*1000 {
		suffix = "PB"
		total /= 1000 * 1000
	} else if total > 1000 {
		suffix = "TB"
		total /= 1000
	}

	typ := "HDD"
	if d.NVMe {
		typ = "NVMe"
	} else if d.SSD {
		typ = "SSD"
	}

	return fmt.Sprintf("%d%s-%s", total, suffix, typ)
}

var diskRE = regexp.MustCompile(`(?i)(?:(\d+) x )?(\d+)(?:GB| GB)?( NVMe)?(?: (SSD|HDD))?`)

func parseStorage(s string) (Disk, error) {
	var d Disk
	if s == "EBS only" {
		return d, nil
	}

	m := diskRE.FindStringSubmatch(s)

	if len(m) < 1 {
		return d, fmt.Errorf("parse storage fail for %q", s)
	}

	d.Count = 1

	if m[1] != "" {
		d.Count, _ = strconv.Atoi(m[1])
	}

	d.PerDiskGB, _ = strconv.Atoi(m[2])

	if m[3] != "" {
		d.NVMe = true
		d.SSD = true
	}

	// If type is specified as SSD or if GB is present without type (i8g's do this)
	if m[4] == "SSD" || (m[4] == "" && strings.Contains(s, "GB")) {
		d.SSD = true
	}

	return d, nil
}
