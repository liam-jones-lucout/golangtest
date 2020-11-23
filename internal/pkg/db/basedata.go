package db

import "github.com/liam-jones-lucout/golangtest/internal/pkg/spaceshipmodels"

func GetBaseData() spaceshipmodels.Spaceships {
	return spaceshipmodels.Spaceships{
		{
			Name:   "Scarier triangle",
			Class:  "Executor",
			Crew:   65146,
			Image:  "https://fillmurray.com/200/300",
			Status: "operational",
			Value:  199999,
			Armament: []*spaceshipmodels.Armament{
				{
					Title: "Blaster",
					Qty:   42,
				},
				{
					Title: "OtherBlaster",
					Qty:   416,
				},
				{
					Title: "Shield Buster",
					Qty:   1,
				},
			},
		},
		{
			Name:   "Ackbars trap",
			Class:  "Home One",
			Crew:   65146,
			Image:  "https://fillmurray.com/200/300",
			Status: "operational",
			Value:  199999,
			Armament: []*spaceshipmodels.Armament{
				{
					Title: "Blaster",
					Qty:   42,
				},
				{
					Title: "OtherBlaster",
					Qty:   416,
				},
				{
					Title: "Shield Buster",
					Qty:   1,
				},
			},
		},
		{
			Name:   "Ackbars Other Ride",
			Class:  "Home One",
			Crew:   65146,
			Image:  "https://fillmurray.com/200/300",
			Status: "damaged",
			Value:  199999,
			Armament: []*spaceshipmodels.Armament{
				{
					Title: "Blaster",
					Qty:   42,
				},
				{
					Title: "OtherBlaster",
					Qty:   416,
				},
				{
					Title: "Shield Buster",
					Qty:   1,
				},
			},
		},
		{
			Name:   "Late Arrival",
			Class:  "Imperial landing craft",
			Crew:   65146,
			Image:  "https://fillmurray.com/200/300",
			Status: "operational",
			Value:  199999,
			Armament: []*spaceshipmodels.Armament{
				{
					Title: "Blaster",
					Qty:   42,
				},
				{
					Title: "OtherBlaster",
					Qty:   416,
				},
				{
					Title: "Shield Buster",
					Qty:   1,
				},
			},
		},
		{
			Name:   "Shootybob",
			Class:  "Imperial shuttle",
			Crew:   65146,
			Image:  "https://fillmurray.com/200/300",
			Status: "operational",
			Value:  199999,
			Armament: []*spaceshipmodels.Armament{
				{
					Title: "Blaster",
					Qty:   42,
				},
				{
					Title: "OtherBlaster",
					Qty:   416,
				},
				{
					Title: "Shield Buster",
					Qty:   1,
				},
			},
		},
		{
			Name:   "Scary triangle",
			Class:  "Imperial Star Destroyer",
			Crew:   65146,
			Image:  "https://fillmurray.com/200/300",
			Status: "operational",
			Value:  199999,
			Armament: []*spaceshipmodels.Armament{
				{
					Title: "Blaster",
					Qty:   42,
				},
				{
					Title: "OtherBlaster",
					Qty:   416,
				},
				{
					Title: "Shield Buster",
					Qty:   1,
				},
			},
		},
		{
			Name:   "Millennium Falcon",
			Class:  "YT-1300 light freighter",
			Crew:   65146,
			Image:  "https://fillmurray.com/200/300",
			Status: "operational",
			Value:  199999,
			Armament: []*spaceshipmodels.Armament{
				{
					Title: "Blaster",
					Qty:   42,
				},
				{
					Title: "OtherBlaster",
					Qty:   416,
				},
				{
					Title: "Shield Buster",
					Qty:   1,
				},
			},
		},
		{
			Name:   "Curie",
			Class:  "Rebel Medical Frigate",
			Crew:   65146,
			Image:  "https://fillmurray.com/200/300",
			Status: "operational",
			Value:  199999,
			Armament: []*spaceshipmodels.Armament{
				{
					Title: "Blaster",
					Qty:   42,
				},
				{
					Title: "OtherBlaster",
					Qty:   416,
				},
				{
					Title: "Shield Buster",
					Qty:   1,
				},
			},
		},
		{
			Name:   "Zippy",
			Class:  "Rebel Transport",
			Crew:   65146,
			Image:  "https://fillmurray.com/200/300",
			Status: "operational",
			Value:  199999,
			Armament: []*spaceshipmodels.Armament{
				{
					Title: "Blaster",
					Qty:   42,
				},
				{
					Title: "OtherBlaster",
					Qty:   416,
				},
				{
					Title: "Shield Buster",
					Qty:   1,
				},
			},
		},
	}
}
