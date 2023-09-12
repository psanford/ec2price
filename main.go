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

		disk, err := parseStorage(attrs.Storage)
		if err != nil {
			log.Printf("parse storage err: %s", err)
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
		Flags:  GpuNvidiaSuffix,
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
		Flags:  GpuAmdSuffix,
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
	{
		// https://aws.amazon.com/blogs/aws/new-amazon-ec2-x2gd-instances-graviton2-power-for-memory-intensive-workloads/
		Name:   "x2gd",
		Year:   2021,
		Prefix: MemXtremePrefix,
		Flags:  GravitonSuffix | NVMeSuffix,
	},
	{
		// https://aws.amazon.com/blogs/aws/new-amazon-ec2-m6i-instances-powered-by-the-latest-generation-intel-xeon-scalable-processors/
		Name:   "m6i",
		Year:   2022,
		Prefix: MainPrefix,
		Flags:  IntelSuffix,
	},
	{
		// https://aws.amazon.com/about-aws/whats-new/2021/09/amazon-ec2-vt1-instances-video-transcoding/
		Name:   "vt1",
		Year:   2021,
		Prefix: VideoTranscodingPrefix,
	},
	{
		// https://aws.amazon.com/blogs/aws/new-amazon-ec2-c6i-instances-powered-by-the-latest-generation-intel-xeon-scalable-processors/
		Name:   "c6i",
		Year:   2022,
		Prefix: CpuPrefix,
		Flags:  IntelSuffix,
	},
	{
		// https://aws.amazon.com/blogs/aws/new-ec2-instances-g5-with-nvidia-a10g-tensor-core-gpus/
		Name:   "g5",
		Year:   2021,
		Prefix: GPUPrefix,
	},
	{
		// https://aws.amazon.com/blogs/aws/new-amazon-ec2-r6i-memory-optimized-instances-powered-by-the-latest-generation-intel-xeon-scalable-processors/
		Name:   "r6i",
		Year:   2021,
		Prefix: MemMorePrefix,
		Flags:  IntelSuffix,
	},
	{
		// https://aws.amazon.com/blogs/aws/new-amazon-ec2-m6a-instances-powered-by-3rd-gen-amd-epyc-processors/
		Name:   "m6a",
		Year:   2021,
		Prefix: MainPrefix,
		Flags:  AmdSuffix,
	},
	{
		// https://aws.amazon.com/blogs/aws/new-amazon-ec2-g5g-instances-powered-by-aws-graviton2-processors-and-nvidia-t4g-tensor-core-gpus/
		Name:   "g5g",
		Year:   2021,
		Prefix: GPUPrefix,
		Flags:  GravitonSuffix,
	},
	{
		// https://aws.amazon.com/blogs/aws/join-the-preview-amazon-ec2-c7g-instances-powered-by-new-aws-graviton3-processors/
		// graviton 3
		Name:   "c7g",
		Year:   2021,
		Prefix: CpuPrefix,
		Flags:  GravitonSuffix,
	},
	{
		// https://aws.amazon.com/blogs/aws/new-storage-optimized-amazon-ec2-instances-im4gn-and-is4gen-powered-by-aws-graviton2-processors/
		Name:   "im4gn",
		Year:   2021,
		Prefix: SSDPrefix,
		Flags:  GravitonSuffix | NVMeSuffix,
	},
	{
		// https://aws.amazon.com/blogs/aws/new-storage-optimized-amazon-ec2-instances-im4gn-and-is4gen-powered-by-aws-graviton2-processors/
		Name:   "is4gn",
		Year:   2021,
		Prefix: SSDPrefix,
		Flags:  GravitonSuffix | NVMeSuffix,
	},
	{
		// https://aws.amazon.com/about-aws/whats-new/2021/11/amazon-ec2-trn1-instances/
		// trn:"Trainium"
		Name:   "trn1",
		Year:   2021,
		Prefix: InferencePrefix,
	},
	{
		// https://aws.amazon.com/blogs/aws/new-amazon-ec2-hpc6a-instance-optimized-for-high-performance-computing/
		Name:   "hpc6a",
		Year:   2022,
		Prefix: HPCPrefix,
		Flags:  AmdSuffix,
	},
	{
		// https://aws.amazon.com/blogs/aws/new-amazon-ec2-x2iezn-instances-powered-by-the-fastest-intel-xeon-scalable-cpu-for-memory-intensive-workloads/
		Name:   "x2iezn",
		Year:   2022,
		Prefix: XeonScalablePrefix,
		Flags:  IntelSuffix | ExtendedMemorySuffix | HighFreqSuffix | NVMeSuffix,
	},
	{
		// https://aws.amazon.com/blogs/aws/new-amazon-ec2-c6a-instances-powered-by-3rd-gen-amd-epyc-processors-for-compute-intensive-workloads/
		Name:   "c6a",
		Year:   2022,
		Prefix: CpuPrefix,
		Flags:  AmdSuffix,
	},
	{
		// https://aws.amazon.com/blogs/aws/new-amazon-ec2-x2idn-and-x2iedn-instances-for-memory-intensive-workloads-with-higher-network-bandwidth/
		Name:   "x2iedn",
		Year:   2022,
		Prefix: XeonScalablePrefix,
		Flags:  IntelSuffix | ExtendedMemorySuffix | NetworkSuffix | NVMeSuffix,
	},
	{
		// https://aws.amazon.com/blogs/aws/new-amazon-ec2-x2idn-and-x2iedn-instances-for-memory-intensive-workloads-with-higher-network-bandwidth/
		Name:   "x2idn",
		Year:   2022,
		Prefix: XeonScalablePrefix,
		Flags:  IntelSuffix | NetworkSuffix | NVMeSuffix,
	},
	{
		// https://aws.amazon.com/blogs/aws/new-storage-optimized-amazon-ec2-instances-i4i-powered-by-intel-xeon-scalable-ice-lake-processors/
		Name:   "i4i",
		Year:   2022,
		Prefix: SSDPrefix,
		Flags:  IntelSuffix,
	},
	{
		// https://aws.amazon.com/about-aws/whats-new/2022/05/amazon-ec2-p4de-gpu-instances-ml-training-hpc/
		Name:   "p4de",
		Year:   2022,
		Prefix: GPUPrefix,
		Flags:  NVMeSuffix | ExtendedMemorySuffix,
	},
	{
		// https://aws.amazon.com/blogs/aws/new-amazon-ec2-m6id-and-c6id-instances-with-up-to-7-6-tb-local-nvme-storage/
		Name:   "c6id",
		Year:   2022,
		Prefix: CpuPrefix,
		Flags:  IntelSuffix | NVMeSuffix,
	},
	{
		// https://aws.amazon.com/blogs/aws/new-amazon-ec2-m6id-and-c6id-instances-with-up-to-7-6-tb-local-nvme-storage/
		Name:   "m6id",
		Year:   2022,
		Prefix: MainPrefix,
		Flags:  IntelSuffix | NVMeSuffix,
	},
	{
		// https://aws.amazon.com/blogs/aws/new-amazon-ec2-r6id-instances/
		Name:   "r6id",
		Year:   2022,
		Prefix: MemMorePrefix,
		Flags:  IntelSuffix | NVMeSuffix,
	},
	{
		// https://aws.amazon.com/blogs/aws/new-amazon-ec2-r6a-instances-powered-by-3rd-gen-amd-epyc-processors-for-memory-intensive-workloads/
		Name:   "r6a",
		Year:   2022,
		Prefix: MemMorePrefix,
		Flags:  AmdSuffix,
	},
	{
		// https://aws.amazon.com/about-aws/whats-new/2022/11/introducing-amazon-ec2-r7iz-instances/
		// https://aws.amazon.com/blogs/aws/new-amazon-ec2-instance-types-in-the-works-c7gn-r7iz-and-hpc7g/
		Name:   "r7iz",
		Year:   2022,
		Prefix: MemMorePrefix,
		Flags:  IntelSuffix | HighFreqSuffix,
	},
	{
		// https://aws.amazon.com/blogs/aws/new-amazon-ec2-m6id-and-c6id-instances-with-up-to-7-6-tb-local-nvme-storage/
		Name:   "m6in",
		Year:   2022,
		Prefix: MainPrefix,
		Flags:  IntelSuffix | NetworkSuffix,
	},
	{
		// https://aws.amazon.com/blogs/aws/new-amazon-ec2-m6id-and-c6id-instances-with-up-to-7-6-tb-local-nvme-storage/
		Name:   "m6idn",
		Year:   2022,
		Prefix: MainPrefix,
		Flags:  IntelSuffix | NVMeSuffix | NetworkSuffix,
	},
	{
		// https://aws.amazon.com/blogs/aws/new-amazon-ec2-m6id-and-c6id-instances-with-up-to-7-6-tb-local-nvme-storage/
		Name:   "c6in",
		Year:   2022,
		Prefix: CpuPrefix,
		Flags:  IntelSuffix | NetworkSuffix,
	},
	{
		// https://aws.amazon.com/blogs/aws/new-amazon-ec2-m6id-and-c6id-instances-with-up-to-7-6-tb-local-nvme-storage/
		Name:   "r6in",
		Year:   2022,
		Prefix: MemMorePrefix,
		Flags:  IntelSuffix | NetworkSuffix,
	},
	{
		// https://aws.amazon.com/blogs/aws/new-amazon-ec2-m6id-and-c6id-instances-with-up-to-7-6-tb-local-nvme-storage/
		Name:   "r6idn",
		Year:   2022,
		Prefix: MemMorePrefix,
		Flags:  IntelSuffix | NVMeSuffix | NetworkSuffix,
	},
	{
		// https://aws.amazon.com/blogs/aws/new-amazon-ec2-instance-types-in-the-works-c7gn-r7iz-and-hpc7g/
		Name:   "c7gn",
		Year:   2022,
		Prefix: CpuPrefix,
		Flags:  GravitonSuffix | NetworkSuffix,
	},
	{
		// https://aws.amazon.com/blogs/aws/new-amazon-ec2-instance-types-in-the-works-c7gn-r7iz-and-hpc7g/
		Name:   "hpc7g",
		Year:   2022,
		Prefix: HPCPrefix,
		Flags:  GravitonSuffix,
	},
	{
		// https://aws.amazon.com/blogs/aws/new-amazon-ec2-hpc6id-instances-optimized-for-high-performance-computing/
		Name:   "hpc6id",
		Year:   2022,
		Prefix: HPCPrefix,
		Flags:  IntelSuffix | NVMeSuffix,
	},
	{
		// https://aws.amazon.com/blogs/aws/new-graviton3-based-general-purpose-m7g-and-memory-optimized-r7g-amazon-ec2-instances
		Name:   "m7g",
		Year:   2023,
		Prefix: MainPrefix,
		Flags:  GravitonSuffix,
	},
	{
		// https://aws.amazon.com/blogs/aws/new-graviton3-based-general-purpose-m7g-and-memory-optimized-r7g-amazon-ec2-instances
		Name:   "r7g",
		Year:   2023,
		Prefix: MemMorePrefix,
		Flags:  GravitonSuffix,
	},
	{
		// https://aws.amazon.com/blogs/aws/amazon-ec2-inf2-instances-for-low-cost-high-performance-generative-ai-inference-are-now-generally-available/
		Name:   "inf2",
		Year:   2023,
		Prefix: InferencePrefix,
	},
	{
		// https://aws.amazon.com/blogs/aws/new-storage-optimized-amazon-ec2-i4g-instances-graviton-processors-and-aws-nitro-ssds/
		Name:   "i4g",
		Year:   2023,
		Prefix: SSDPrefix,
		Flags:  GravitonSuffix,
	},
	{
		// https://aws.amazon.com/blogs/aws/new-amazon-ec2-p5-instances-powered-by-nvidia-h100-tensor-core-gpus-for-accelerating-generative-ai-and-hpc-applications/
		Name:   "p5",
		Year:   2023,
		Prefix: GPUPrefix,
	},
	{
		// https://aws.amazon.com/blogs/aws/new-amazon-ec2-instances-c7gd-m7gd-and-r7gd-powered-by-aws-graviton3-processor-with-local-nvme-based-ssd-storage/
		Name:   "c7gd",
		Year:   2023,
		Prefix: CpuPrefix,
		Flags:  GravitonSuffix | NVMeSuffix,
	},
	{
		// https://aws.amazon.com/blogs/aws/new-amazon-ec2-instances-c7gd-m7gd-and-r7gd-powered-by-aws-graviton3-processor-with-local-nvme-based-ssd-storage/
		Name:   "m7gd",
		Year:   2023,
		Prefix: MainPrefix,
		Flags:  GravitonSuffix | NVMeSuffix,
	},
	{
		// https://aws.amazon.com/blogs/aws/new-amazon-ec2-instances-c7gd-m7gd-and-r7gd-powered-by-aws-graviton3-processor-with-local-nvme-based-ssd-storage/
		Name:   "r7gd",
		Year:   2023,
		Prefix: MemMorePrefix,
		Flags:  GravitonSuffix | NVMeSuffix,
	},
	{
		// https://aws.amazon.com/blogs/aws/new-seventh-generation-general-purpose-amazon-ec2-instances-m7i-flex-and-m7i/
		Name:   "m7i",
		Year:   2023,
		Prefix: MainPrefix,
		Flags:  IntelSuffix,
	},
	{
		// https://aws.amazon.com/blogs/aws/new-seventh-generation-general-purpose-amazon-ec2-instances-m7i-flex-and-m7i/
		Name:   "m7i-flex",
		Year:   2023,
		Prefix: MainPrefix,
		Flags:  IntelSuffix | FlexSuffix,
	},
	{
		// https://aws.amazon.com/blogs/aws/new-amazon-ec2-m7a-general-purpose-instances-powered-by-4th-gen-amd-epyc-processors/
		Name:   "m7a",
		Year:   2023,
		Prefix: MainPrefix,
		Flags:  AmdSuffix,
	},
	{
		// https://aws.amazon.com/blogs/aws/new-amazon-ec2-hpc7a-instances-powered-by-4th-gen-amd-epyc-processors-optimized-for-high-performance-computing/
		Name:   "hpc7a",
		Year:   2023,
		Prefix: HPCPrefix,
		Flags:  AmdSuffix,
	},
	{
		// https://aws.amazon.com/blogs/aws/new-amazon-ec2-r7a-instances-powered-by-4th-gen-amd-epyc-processors-for-memory-optimized-workloads/
		Name:   "r7a",
		Year:   2023,
		Prefix: MemMorePrefix,
		Flags:  AmdSuffix,
	},
}

type InstanceCodePrefix int

const (
	ArmPrefix              InstanceCodePrefix = iota // a
	BurstPrefix                                      // t
	MainPrefix                                       // m
	CpuPrefix                                        // c
	MemMorePrefix                                    // r
	MemXtremePrefix                                  // x
	HighFreqPrefix                                   // z
	GPUPrefix                                        // p
	InferencePrefix                                  // inf,trn("Trainium")
	FPGAPrefix                                       // f
	SSDPrefix                                        // i
	DenseHDDPrefix                                   // d
	ClusterComputePrefix                             // cc,cr
	VideoTranscodingPrefix                           // vt
	HPCPrefix                                        // hpc
	XeonScalablePrefix                               // x2
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

var diskRE = regexp.MustCompile(`(?:(\d+) x )?(\d+)(?: GB)?( NVMe)? (SSD|HDD)`)

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

	if m[4] == "SSD" {
		d.SSD = true
	}

	return d, nil
}
