package pnr

// func convertRemarks(res AlaskaManagePnrResponse, pnr *PNR) {
// 	for _, collection := range res.Collection {
// 		for _, remark := range collection.Remarks {
// 			pnr.Remarks = append(pnr.Remarks, Remark{
// 				FreeFormText: remark,
// 				RemarkType:   "RRMK",
// 			})
// 		}
// 	}
// }

func convertItinerary(res AlaskaManagePnrResponse, pnr *PNR) {
	pnr.Itinerary.Origin = res.Itinerary.Origin
	pnr.Itinerary.Type = res.Itinerary.TripType
	pnr.Itinerary.MatchesTickets = res.Itinerary.MatchesTickets

}

func convertFlights(res AlaskaManagePnrResponse, pnr *PNR) {
	for _, flight := range res.Itinerary.ItinerarySlices {
		for _, segment := range flight.Segments {
			f := Flight{
				OriginAirportCode:      segment.DepartureAirport,
				DestinationAirportCode: segment.ArrivalAirport,
				OperatingAirlineCode:   segment.OperatingAirlineCode,
				MarketingAirlineCode:   segment.OperatingAirlineCode,
				FlightNumber:           segment.OperatingFlightNumber,
				CurrentActionCode:      segment.ActionCode,
				ClassOfService:         segment.ClassOfService,
				Cabin:                  segment.Cabin,
				ScheduledDeparture:     segment.ScheduledDepartureDateTime,
				ScheduledArrival:       segment.ScheduledArrivalDateTime,
				IsDisrupted:            segment.IsDisrupted,
				IsFlown:                segment.IsFlown,
				Distance:               segment.Distance,
			}
			for _, ssr := range segment.SpecialServiceRequests {
				f.SSRs = append(f.SSRs, SSR{
					Remark:      ssr.FreeText,
					Type:        ssr.ServiceCode,
					AirlineCode: ssr.VendorCode,
					Status:      ssr.ActionCode,
					FlightNum:   ssr.FlightNumber,
					Id:          ssr.ID,
				})
			}
			pnr.Flights = append(pnr.Flights, f)
		}
	}
}

func convertPassengers(res AlaskaManagePnrResponse, pnr *PNR) {
	for _, pax := range res.Passengers {
		p := Passenger{
			Name:                   pax.FirstName + " " + pax.LastName,
			OverbookingEligible:    false,
			TierStatus:             pax.TierStatus,
			FfNumber:               pax.LoyaltyNumber,
			FeeForSeatSelection:    true,
			FeeForAmPlusUpgrade:    true,
			FeeForPreferredUpgrade: true,
		}
		for _, ssr := range pax.SpecialServiceRequests {

			p.SSRs = append(p.SSRs, SSR{
				Remark:      ssr.FreeText,
				Type:        ssr.ServiceCode,
				AirlineCode: ssr.VendorCode,
				Status:      ssr.ActionCode,
				FlightNum:   ssr.FlightNumber,
				FlightDate:  ssr.FlightDate,
			})
		}
		pnr.Passengers = append(pnr.Passengers, p)
	}
}
func convertTickets(res AlaskaManagePnrResponse, pnr *PNR) {
	for _, pax := range res.Passengers {
		for _, ticket := range pax.Tickets {
			t := Ticket{
				Number: ticket.Number,
				// ExpirationDate:         ticket.ExpirationDate,
				// IssueDate:              ticket.IssueDate,
				Type:          ticket.Type,
				Designator:    ticket.Designator,
				Status:        ticket.IsActive,
				PassengerName: pax.FirstName + " " + pax.LastName,
				NumCoupons:    uint64(len(ticket.Coupons)),
				// ValidatedAgainstCoupon: couponsMatchFlights(res, ticket.Number),
			}
			for _, coupon := range ticket.Coupons {
				t.Coupons = append(t.Coupons, Coupon{
					Index:       coupon.Index,
					Origin:      coupon.Origin,
					Destination: coupon.Destination,
					Status:      coupon.Status,
				})
			}

			pnr.Tickets = append(pnr.Tickets, t)
		}
	}
}

func convertEarnings(_ AlaskaManagePnrResponse, pnr *PNR) {
	pnr.SMCalcLink = generateSmcalcLink(pnr)
}
