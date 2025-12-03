package main

import (
	"reflect"
	"testing"
)

func TestParseNetworkPerf(t *testing.T) {
	type TC struct {
		in  string
		out NetworkPerf
	}

	cases := []TC{
		{
			in: "100000 Megabit",
			out: NetworkPerf{
				CapGb: 100.000,
			},
		},
		{
			in: "100 Gigabit",
			out: NetworkPerf{
				CapGb: 100,
			},
		},
		{
			in: "10 Gigabit",
			out: NetworkPerf{
				CapGb: 10,
			},
		},
		{
			in: "12500 Megabit",
			out: NetworkPerf{
				CapGb: 12.500,
			},
		},
		{
			in: "12 Gigabit",
			out: NetworkPerf{
				CapGb: 12,
			},
		},
		{
			in: "150000 Megabit",
			out: NetworkPerf{
				CapGb: 150.000,
			},
		},
		{
			in: "150 Gigabit",
			out: NetworkPerf{
				CapGb: 150,
			},
		},
		{
			in: "15 Gigabit",
			out: NetworkPerf{
				CapGb: 15,
			},
		},
		{
			in: "1600 Gigabit",
			out: NetworkPerf{
				CapGb: 1600,
			},
		},
		{
			in: "18750 Megabit",
			out: NetworkPerf{
				CapGb: 18.750,
			},
		},
		{
			in: "200000 Megabit",
			out: NetworkPerf{
				CapGb: 200.000,
			},
		},
		{
			in: "200 Gigabit",
			out: NetworkPerf{
				CapGb: 200,
			},
		},
		{
			in: "20 Gigabit",
			out: NetworkPerf{
				CapGb: 20,
			},
		},
		{
			in: "22500 Megabit",
			out: NetworkPerf{
				CapGb: 22.500,
			},
		},
		{
			in: "25000 Megabit",
			out: NetworkPerf{
				CapGb: 25.000,
			},
		},
		{
			in: "25 Gigabit",
			out: NetworkPerf{
				CapGb: 25,
			},
		},
		{
			in: "30 Gigabit",
			out: NetworkPerf{
				CapGb: 30,
			},
		},
		{
			in: "3125 Megabit",
			out: NetworkPerf{
				CapGb: 3.125,
			},
		},
		{
			in: "37500 Megabit",
			out: NetworkPerf{
				CapGb: 37.500,
			},
		},
		{
			in: "400 Gigabit",
			out: NetworkPerf{
				CapGb: 400,
			},
		},
		{
			in: "40 Gigabit",
			out: NetworkPerf{
				CapGb: 40,
			},
		},
		{
			in: "50000 Megabit",
			out: NetworkPerf{
				CapGb: 50.000,
			},
		},
		{
			in: "50 Gigabit",
			out: NetworkPerf{
				CapGb: 50,
			},
		},
		{
			in: "6250 Megabit",
			out: NetworkPerf{
				CapGb: 6.250,
			},
		},
		{
			in: "75000 Megabit",
			out: NetworkPerf{
				CapGb: 75.000,
			},
		},
		{
			in: "75 Gigabit",
			out: NetworkPerf{
				CapGb: 75,
			},
		},
		{
			in: "800 Gigabit",
			out: NetworkPerf{
				CapGb: 800,
			},
		},
		{
			in: "Up to 10 Gigabit",
			out: NetworkPerf{
				CapGb:    10,
				Bursting: true,
			},
		},
		{
			in: "Up to 12500 Megabit",
			out: NetworkPerf{
				CapGb:    12.500,
				Bursting: true,
			},
		},
		{
			in: "Up to 12 Gigabit",
			out: NetworkPerf{
				CapGb:    12,
				Bursting: true,
			},
		},
		{
			in: "Up to 15 Gigabit",
			out: NetworkPerf{
				CapGb:    15,
				Bursting: true,
			},
		},
		{
			in: "Up to 25000 Megabit",
			out: NetworkPerf{
				CapGb:    25.000,
				Bursting: true,
			},
		},
		{
			in: "Up to 25 Gigabit",
			out: NetworkPerf{
				CapGb:    25,
				Bursting: true,
			},
		},
		{
			in: "Up to 30000 Megabit",
			out: NetworkPerf{
				CapGb:    30.000,
				Bursting: true,
			},
		},
		{
			in: "Up to 30 Gigabit",
			out: NetworkPerf{
				CapGb:    30,
				Bursting: true,
			},
		},
		{
			in: "Up to 40000 Megabit",
			out: NetworkPerf{
				CapGb:    40.000,
				Bursting: true,
			},
		},
		{
			in: "Up to 40 Gigabit",
			out: NetworkPerf{
				CapGb:    40,
				Bursting: true,
			},
		},
		{
			in: "Up to 50000 Megabit",
			out: NetworkPerf{
				CapGb:    50.000,
				Bursting: true,
			},
		},
		{
			in: "Up to 50 Gigabit",
			out: NetworkPerf{
				CapGb:    50,
				Bursting: true,
			},
		},
		{
			in: "Up to 5 Gigabit",
			out: NetworkPerf{
				CapGb:    5,
				Bursting: true,
			},
		},

		{
			in: "Very Low",
			out: NetworkPerf{
				CapGb:    0.01,
				Bursting: true,
			},
		},
		{
			in: "High",
			out: NetworkPerf{
				CapGb: 1,
			},
		},
		{
			in: "Low",
			out: NetworkPerf{
				CapGb:    0.01,
				Bursting: true,
			},
		},
		{
			in: "Low to Moderate",
			out: NetworkPerf{
				CapGb:    0.01,
				Bursting: true,
			},
		},
		{
			in: "Moderate",
			out: NetworkPerf{
				CapGb:    0.1,
				Bursting: true,
			},
		},
		{
			in: "NA",
			out: NetworkPerf{
				CapGb: 1,
			},
		},
	}

	for _, tc := range cases {
		got, err := parseNetPerf(tc.in)
		if err != nil {
			t.Errorf("parse %q err: %s", tc.in, err)
		}
		if !reflect.DeepEqual(got, tc.out) {
			t.Errorf("%q parse mismatch: got=%+v exp=%+v", tc.in, got, tc.out)
		}
	}
}

func TestParseStorage(t *testing.T) {
	type TC struct {
		in  string
		out Disk
	}

	cases := []TC{
		{
			in: "1200 GB NVMe SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 1200,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "125 GB NVMe SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 125,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "12 x 14000 HDD",
			out: Disk{
				Count:     12,
				PerDiskGB: 14000,
			},
		},
		{
			in: "12 x 2000 HDD",
			out: Disk{
				Count:     12,
				PerDiskGB: 2000,
			},
		},
		{
			in: "150 GB NVMe SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 150,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "16 x 14000 HDD",
			out: Disk{
				Count:     16,
				PerDiskGB: 14000,
			},
		},
		{
			in: "1 x 100 NVMe SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 100,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "1 x 118 NVMe SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 118,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "1 x 118 SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 118,
				SSD:       true,
			},
		},
		{
			in: "1 x 120 SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 120,
				SSD:       true,
			},
		},
		{
			in: "1 x 1250 NVMe SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 1250,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "1 x 150 NVMe SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 150,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "1 x 160 SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 160,
				SSD:       true,
			},
		},
		{
			in: "1 x 1875 NVMe SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 1875,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "1 x 1875 SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 1875,
				SSD:       true,
			},
		},
		{
			in: "1 x 1900 GB NVMe SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 1900,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "1 x 1900 NVMe SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 1900,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "1 x 1900 SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 1900,
				SSD:       true,
			},
		},
		{
			in: "1 x 1920 SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 1920,
				SSD:       true,
			},
		},
		{
			in: "1 x 2000 HDD",
			out: Disk{
				Count:     1,
				PerDiskGB: 2000,
			},
		},
		{
			in: "1 x 200 NVMe SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 200,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "1 x 237 NVMe SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 237,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "1 x 237 SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 237,
				SSD:       true,
			},
		},
		{
			in: "1 x 240 SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 240,
				SSD:       true,
			},
		},
		{
			in: "1 x 2500 NVMe SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 2500,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "1 x 250 GB NVMe SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 250,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "1 x 300 NVMe SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 300,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "1 x 320 SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 320,
				SSD:       true,
			},
		},
		{
			in: "1 x 32 SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 32,
				SSD:       true,
			},
		},
		{
			in: "1 x 350 SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 350,
				SSD:       true,
			},
		},
		{
			in: "1 x 3750 NVMe SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 3750,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "1 x 3750 SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 3750,
				SSD:       true,
			},
		},
		{
			in: "1 x 3800 GB NVMe SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 3800,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "1 x 400 NVMe SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 400,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "1 x 410 SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 410,
				SSD:       true,
			},
		},
		{
			in: "1 x 420 SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 420,
				SSD:       true,
			},
		},
		{
			in: "1 x 450 GB NVMe SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 450,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "1 x 450 NVMe SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 450,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "1 x 468 NVMe SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 468,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "1 x 468 SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 468,
				SSD:       true,
			},
		},
		{
			in: "1 x 470 NVMe SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 470,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "1 x 474 SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 474,
				SSD:       true,
			},
		},
		{
			in: "1 x 475 NVMe SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 475,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "1 x 480 SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 480,
				SSD:       true,
			},
		},
		{
			in: "1 x 4 SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 4,
				SSD:       true,
			},
		},
		{
			in: "1 x 50 NVMe SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 50,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "1 x 59 NVMe SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 59,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "1 x 59 SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 59,
				SSD:       true,
			},
		},
		{
			in: "1 x 600 GB NVMe SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 600,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "1 x 60 SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 60,
				SSD:       true,
			},
		},
		{
			in: "1 x 7500 NVMe SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 7500,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "1 x 7500 SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 7500,
				SSD:       true,
			},
		},
		{
			in: "1 x 75 NVMe SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 75,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "1 x 800 SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 800,
				SSD:       true,
			},
		},
		{
			in: "1 x 80 SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 80,
				SSD:       true,
			},
		},
		{
			in: "1 x 850 SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 850,
				SSD:       true,
			},
		},
		{
			in: "1 x 900 GB NVMe SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 900,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "1 x 900 NVMe SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 900,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "1 x 937 NVMe SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 937,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "1 x 937 SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 937,
				SSD:       true,
			},
		},
		{
			in: "1 x 940 NVMe SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 940,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "1 x 950 NVMe SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 950,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "1 x 950 SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 950,
				SSD:       true,
			},
		},
		{
			in: "1 x 960 SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 960,
				SSD:       true,
			},
		},
		{
			in: "225 GB NVMe SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 225,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "2400 GB NVMe SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 2400,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "24 x 14000 HDD",
			out: Disk{
				Count:     24,
				PerDiskGB: 14000,
			},
		},
		{
			in: "24 x 2000 HDD",
			out: Disk{
				Count:     24,
				PerDiskGB: 2000,
			},
		},
		{
			in: "2 x 1200 NVMe SSD",
			out: Disk{
				Count:     2,
				PerDiskGB: 1200,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "2 x 120 SSD",
			out: Disk{
				Count:     2,
				PerDiskGB: 120,
				SSD:       true,
			},
		},
		{
			in: "2 x 14000 HDD",
			out: Disk{
				Count:     2,
				PerDiskGB: 14000,
			},
		},
		{
			in: "2 x 1425 NVMe SSD",
			out: Disk{
				Count:     2,
				PerDiskGB: 1425,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "2 x 1425 SSD",
			out: Disk{
				Count:     2,
				PerDiskGB: 1425,
				SSD:       true,
			},
		},
		{
			in: "2 x 160 SSD",
			out: Disk{
				Count:     2,
				PerDiskGB: 160,
				SSD:       true,
			},
		},
		{
			in: "2 x 16 SSD",
			out: Disk{
				Count:     2,
				PerDiskGB: 16,
				SSD:       true,
			},
		},
		{
			in: "2 x 1900 NVMe SSD",
			out: Disk{
				Count:     2,
				PerDiskGB: 1900,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "2 x 1900 SSD",
			out: Disk{
				Count:     2,
				PerDiskGB: 1900,
				SSD:       true,
			},
		},
		{
			in: "2 x 1920 SSD",
			out: Disk{
				Count:     2,
				PerDiskGB: 1920,
				SSD:       true,
			},
		},
		{
			in: "2 x 2000 HDD",
			out: Disk{
				Count:     2,
				PerDiskGB: 2000,
			},
		},
		{
			in: "2 x 2500 NVMe SSD",
			out: Disk{
				Count:     2,
				PerDiskGB: 2500,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "2 x 300 NVMe SSD",
			out: Disk{
				Count:     2,
				PerDiskGB: 300,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "2 x 320 SSD",
			out: Disk{
				Count:     2,
				PerDiskGB: 320,
				SSD:       true,
			},
		},
		{
			in: "2 x 3750 NVMe SSD",
			out: Disk{
				Count:     2,
				PerDiskGB: 3750,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "2 x 3750 SSD",
			out: Disk{
				Count:     2,
				PerDiskGB: 3750,
				SSD:       true,
			},
		},
		{
			in: "2 x 3800 GB NVMe SSD",
			out: Disk{
				Count:     2,
				PerDiskGB: 3800,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "2 x 40 SSD",
			out: Disk{
				Count:     2,
				PerDiskGB: 40,
				SSD:       true,
			},
		},
		{
			in: "2 x 420 SSD",
			out: Disk{
				Count:     2,
				PerDiskGB: 420,
				SSD:       true,
			},
		},
		{
			in: "2 x 600 NVMe SSD",
			out: Disk{
				Count:     2,
				PerDiskGB: 600,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "2 x 7500 NVMe SSD",
			out: Disk{
				Count:     2,
				PerDiskGB: 7500,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "2 x 7500 SSD",
			out: Disk{
				Count:     2,
				PerDiskGB: 7500,
				SSD:       true,
			},
		},
		{
			in: "2 x 800 SSD",
			out: Disk{
				Count:     2,
				PerDiskGB: 800,
				SSD:       true,
			},
		},
		{
			in: "2 x 80 SSD",
			out: Disk{
				Count:     2,
				PerDiskGB: 80,
				SSD:       true,
			},
		},
		{
			in: "2 x 840 SSD",
			out: Disk{
				Count:     2,
				PerDiskGB: 840,
				SSD:       true,
			},
		},
		{
			in: "2 x 900 GB NVMe SSD",
			out: Disk{
				Count:     2,
				PerDiskGB: 900,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "2 x 900 NVMe SSD",
			out: Disk{
				Count:     2,
				PerDiskGB: 900,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "300 GB NVMe SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 300,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "3 x 2000 HDD",
			out: Disk{
				Count:     3,
				PerDiskGB: 2000,
			},
		},
		{
			in: "4 x 1000 GB NVMe SSD",
			out: Disk{
				Count:     4,
				PerDiskGB: 1000,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "4 x 14000 HDD",
			out: Disk{
				Count:     4,
				PerDiskGB: 14000,
			},
		},
		{
			in: "4 x 1425 SSD",
			out: Disk{
				Count:     4,
				PerDiskGB: 1425,
				SSD:       true,
			},
		},
		{
			in: "4 x 1900 NVMe SSD",
			out: Disk{
				Count:     4,
				PerDiskGB: 1900,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "4 x 1900 SSD",
			out: Disk{
				Count:     4,
				PerDiskGB: 1900,
				SSD:       true,
			},
		},
		{
			in: "4 x 2000 HDD",
			out: Disk{
				Count:     4,
				PerDiskGB: 2000,
			},
		},
		{
			in: "4 x 3750 NVMe SSD",
			out: Disk{
				Count:     4,
				PerDiskGB: 3750,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "4 x 3750 SSD",
			out: Disk{
				Count:     4,
				PerDiskGB: 3750,
				SSD:       true,
			},
		},
		{
			in: "4 x 420 SSD",
			out: Disk{
				Count:     4,
				PerDiskGB: 420,
				SSD:       true,
			},
		},
		{
			in: "4 x 600 NVMe SSD",
			out: Disk{
				Count:     4,
				PerDiskGB: 600,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "4 x 7500 NVMe SSD",
			out: Disk{
				Count:     4,
				PerDiskGB: 7500,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "4 x 7500 SSD",
			out: Disk{
				Count:     4,
				PerDiskGB: 7500,
				SSD:       true,
			},
		},
		{
			in: "4 x 800 SSD",
			out: Disk{
				Count:     4,
				PerDiskGB: 800,
				SSD:       true,
			},
		},
		{
			in: "4 x 840 SSD",
			out: Disk{
				Count:     4,
				PerDiskGB: 840,
				SSD:       true,
			},
		},
		{
			in: "4 x 900 NVMe SSD",
			out: Disk{
				Count:     4,
				PerDiskGB: 900,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "4 x 940 NVMe SSD",
			out: Disk{
				Count:     4,
				PerDiskGB: 940,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "600 GB NVMe SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 600,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "6 x 2000 HDD",
			out: Disk{
				Count:     6,
				PerDiskGB: 2000,
			},
		},
		{
			in: "8 x 1000 SSD",
			out: Disk{
				Count:     8,
				PerDiskGB: 1000,
				SSD:       true,
			},
		},
		{
			in: "8 x 14000 HDD",
			out: Disk{
				Count:     8,
				PerDiskGB: 14000,
			},
		},
		{
			in: "8 x 1900 NVMe SSD",
			out: Disk{
				Count:     8,
				PerDiskGB: 1900,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "8 x 2000 HDD",
			out: Disk{
				Count:     8,
				PerDiskGB: 2000,
			},
		},
		{
			in: "8 x 3750 NVMe SSD",
			out: Disk{
				Count:     8,
				PerDiskGB: 3750,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "8 x 3750 SSD",
			out: Disk{
				Count:     8,
				PerDiskGB: 3750,
				SSD:       true,
			},
		},
		{
			in: "8 x 7500 NVMe SSD",
			out: Disk{
				Count:     8,
				PerDiskGB: 7500,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in: "8 x 800 SSD",
			out: Disk{
				Count:     8,
				PerDiskGB: 800,
				SSD:       true,
			},
		},
		{
			in: "900 GB NVMe SSD",
			out: Disk{
				Count:     1,
				PerDiskGB: 900,
				NVMe:      true,
				SSD:       true,
			},
		},
		{
			in:  "EBS only",
			out: Disk{},
		},
		{
			// i8g
			in: "6 x 3750GB",
			out: Disk{
				Count:     6,
				PerDiskGB: 3750,
				SSD:       true,
			},
		},
		{
			in: "1 x 468GB",
			out: Disk{
				Count:     1,
				PerDiskGB: 468,
				SSD:       true,
			},
		},
		{
			in: "2 x 3750GB",
			out: Disk{
				Count:     2,
				PerDiskGB: 3750,
				SSD:       true,
			},
		},
		{
			in: "4 X 940 GB NVMe SSD",
			out: Disk{
				Count:     4,
				PerDiskGB: 940,
				SSD:       true,
				NVMe:      true,
			},
		},
	}

	for _, tc := range cases {
		got, err := parseStorage(tc.in)
		if err != nil {
			t.Errorf("parse %q err: %s", tc.in, err)
		}
		if !reflect.DeepEqual(got, tc.out) {
			t.Errorf("%q parse mismatch: got=%+v exp=%+v", tc.in, got, tc.out)
		}
	}

}

func TestNoDuplicateInstanceTypes(t *testing.T) {
	seen := make(map[string]bool)
	for _, it := range instanceTypes {
		if seen[it.Name] {
			t.Errorf("duplicate instance type: %s", it.Name)
		}
		seen[it.Name] = true
	}
}
