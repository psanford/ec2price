package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

var (
	basePriceURL = "https://pricing.us-east-1.amazonaws.com"
	indexPath    = "/offers/v1.0/aws/index.json"

	region = flag.String("region", "us-east-1", "AWS Region")
)

func main() {
	flag.Parse()

	r, err := http.Get(basePriceURL + indexPath)
	checkErr(err, "Get index")

	if r.StatusCode != 200 {
		b, _ := ioutil.ReadAll(r.Body)
		log.Fatalf("Get index status %d\n%s", r.StatusCode, b)
	}

	dec := json.NewDecoder(r.Body)
	var idx PriceIndex
	err = dec.Decode(&idx)
	checkErr(err, "Read index json")
	r.Body.Close()

	regionURL := idx.Offers["AmazonEC2"].CurrentRegionIndexURL

	r, err = http.Get(basePriceURL + regionURL)
	checkErr(err, "Get region index")

	if r.StatusCode != 200 {
		b, _ := ioutil.ReadAll(r.Body)
		log.Fatalf("Get region index status %d\n%s", r.StatusCode, b)
	}

	dec = json.NewDecoder(r.Body)
	var regionIdx RegionIndex
	err = dec.Decode(&regionIdx)
	checkErr(err, "Read region index json")
	r.Body.Close()

	ec2Path := regionIdx.Regions[*region].CurrentVersionURL

	r, err = http.Get(basePriceURL + ec2Path)
	checkErr(err, "Get prices")

	if r.StatusCode != 200 {
		b, _ := ioutil.ReadAll(r.Body)
		log.Fatalf("Get prices status %d\n%s", r.StatusCode, b)
	}

	dec = json.NewDecoder(r.Body)
	var prices PriceDoc
	err = dec.Decode(&prices)
	checkErr(err, "Read price json")
	r.Body.Close()

	var instances []InstanceType

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
		}

		instances = append(instances, instance)
	}

	sort.Slice(instances, func(a, b int) bool { return instances[a].OnDemandAnnual < instances[b].OnDemandAnnual })

	format := "%15s %10.01f %6s %15s %12d %9.04f %9.02f %.2f\n"
	fmt.Printf("%15s %10s %6s %15s %12s %9s %9s %s\n", "type", "mem", "vcpu", "disk", "disk_total", "hourly", "annual", "annual-reserved")

	for _, in := range instances {
		fmt.Printf(format, in.Name, in.Memory, in.VCPU, in.Disk, in.DiskTotal, in.Hourly, in.OnDemandAnnual, in.ReservedAnnual)
	}

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
