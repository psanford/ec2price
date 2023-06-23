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
