package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"strconv"
	"time"
	"sort"
)

type Restaurant struct {
	RestaurantId string
	Contact string
	A int
	B int
	MinimumReservation int
	MaximumReservation int
	IsStopped	bool	
	IsRemoved bool
}

type DateTiming struct {
	StartTime string
	EndTime string
}

type RestaurantTiming struct {
	RestaurantId string
	Time []string
}

type RestaurantRequestDetail struct {
	RequestId string
	UserId string
	RestaurantId string
	Date int
	Time string
	NoOfPeople int
	IsConfirmed bool
	IsRejected bool
	IsPending bool
	IsCancelled bool
}

type RestaurantTimingDetail struct {
	RestaurantId string
	Time string
}

var Restaurants map[string]Restaurant = make(map[string]Restaurant)
var RestaurantRequestDetails map[string]RestaurantRequestDetail= make(map[string]RestaurantRequestDetail)
var RestaurantTimings map[string]RestaurantTiming = make(map[string]RestaurantTiming)
var RestaurantTimingDetails map[string]RestaurantTimingDetail = make(map[string]RestaurantTimingDetail)
var currentDate int =  int(time.Now().Weekday())
var listCounter int = 0
var commands map[string]string = make(map[string]string)
var isNextDay bool = false

func InitCommands(){
	commands["REQUEST"] = "REQUEST"
	commands["CANCEL"] = "CANCEL"
	commands["CONFIRM"] = "CONFIRM"
	commands["REJECT"] = "REJECT"
	commands["STOP"] = "STOP"
	commands["NEXT_DAY"] = "NEXT_DAY"
	commands["LIST"] = "LIST"
	commands["REMOVE"] = "REMOVE"
}

func NextDay() {
    currentDate = currentDate + 1 
    isNextDay = true
		//tempRestaurantRequestDetails:= RestaurantRequestDetails
		var rejected map[string]string = make(map[string]string)
		

		for key, restaurantRequestDetails := range RestaurantRequestDetails { 
        
				// if (restaurantRequestDetails.IsPending) {
				// 	restaurantRequestDetails.Date+=1
				// }

        if (restaurantRequestDetails.IsPending) && isNextDay {//} (restaurantRequestDetails.Date > currentDate) {
            restaurantRequestDetails.IsRejected = true
						restaurantRequestDetails.IsPending = false
						RestaurantRequestDetails[key] = restaurantRequestDetails
						rejected[restaurantRequestDetails.RequestId] = restaurantRequestDetails.UserId
            //fmt.Printf("to:%v %v has been auto-rejected\n",restaurantRequestDetails.UserId, restaurantRequestDetails.RequestId)
        }

    }

		//sort.

		keys := make([]string, 0, len(rejected))
   	for k := range rejected {
      keys = append(keys, k)
   	}

		sort.Strings(keys)
		for _, k := range keys {
      fmt.Printf("to:%v %v has been auto-rejected\n", rejected[k], k) 
    }

    isNextDay = false
}

func List(){
    listCounter++
}

func IsInsideOfReservationPeriodForStop(restaurantId string, day int, time string) bool {
    counter:= 0
    key:=  restaurantId+"_"+strconv.Itoa(day - 1)
		//fmt.Printf("Inside IsInsideOfReservationPeriodForStop() where key is %s \t %d and value is %s\n", key, day, RestaurantTimingDetails[key].Time)
		_,isPresent := RestaurantTimingDetails[key]
    if !isPresent {
			return false
		}

    restaurantTimingDetails:= RestaurantTimingDetails[key].Time
    //fmt.Println("Inside IsInsideOfReservationPeriodForStop() is :", restaurantTimingDetails)
		var indexPosition int = -1
        indexPosition = strings.Index(restaurantTimingDetails, " ")    
    
    if indexPosition > -1 {
        times:= strings.Split(restaurantTimingDetails, " ")
        
        for i:=0;i<len(times);i++ {
    hours,_:= strconv.Atoi(strings.Split(time, ":")[0])
    tempTime := strings.Split(times[i], "-")[0]
    minHours,_ := strconv.Atoi(strings.Split(tempTime, ":")[0])
    maxHours,_ := strconv.Atoi(strings.Split(tempTime, ":")[1]) 
            if (hours >= minHours) && (hours <= maxHours) {
                counter++;
            }    
        }
        
        if counter == len(times) {
            return true
        }
    } else {
    hours,_:= strconv.Atoi(strings.Split(time, ":")[0])
    tempTime := strings.Split(restaurantTimingDetails, "-")
    minHours,_ := strconv.Atoi(strings.Split(tempTime[0], ":")[0])
    maxHours,_ := strconv.Atoi(strings.Split(tempTime[1], ":")[0])

        if (hours >= minHours) && (hours <= maxHours) {
            return true
        }
    }
    
    return false
}

func Stop(query string){
    restaurantId:= getStringDetails(query,1)
    date,_:= strconv.Atoi(getStringDetails(query,2))
    time:= getStringDetails(query,3)
		restaurant,isPresent := Restaurants[restaurantId]
    
    if !isPresent {
        fmt.Println("Error: No such restaurant")
				return 
    }
    
    if date < currentDate {
        fmt.Println("Error: Specify a date today or after today")
				return 
		}
    
    if IsInsideOfReservationPeriodForStop(restaurantId, date, time ) {
        fmt.Println("Error: Cannot make a reservation already due to being outside the reservation period")
				return 
		}
    
    restaurant.IsStopped = true
		Restaurants[restaurantId] = restaurant
    //fmt.Println("The reservation stop period is set to date entered.")
}

func Reject(query string) {
    //restaurantId:= getStringDetails(query,1)
    reservationId:= getStringDetails(query,2)
		//fmt.Printf("The reservationId is %v\n", reservationId)
    request,isPresent := RestaurantRequestDetails[reservationId]
    
    if !isPresent {
        fmt.Println("Error: No such reservation ID")
        return
    }
    
    if request.IsRejected {
        fmt.Println("Error: Already rejected")
        return
    }
    
    if request.IsConfirmed {
        fmt.Println("Error: Already confirmed")
        return
    }
    
    if request.IsCancelled {
        fmt.Println("Error: Already cancelled")
    }
    
    request.IsRejected = true
    request.IsPending = false
		RestaurantRequestDetails[reservationId] = request
    fmt.Printf("to:%v %v has been rejected\n",request.UserId, reservationId)
}

func IsReservationExistsForUser(userId string, requestId string) bool {
    request := RestaurantRequestDetails[requestId]
    return request.UserId == userId 
}

func IsRemovedRestaurant(restaurantId string) bool{
    return Restaurants[restaurantId].IsRemoved
}

func Cancel(query string) {

    userId := getStringDetails(query,1)
    requestId:= getStringDetails(query,2)
    request := RestaurantRequestDetails[requestId]
    if !IsReservationExistsForUser(userId, requestId) {
        fmt.Println("Error: Not found")
        return 
    }
    
    if IsRemovedRestaurant(request.RestaurantId){
        fmt.Println("Error: No such restaurant")
        return
    }
    
    if request.IsRejected {
        fmt.Println("Error: Rejected")
        return
    }
    
    if request.IsCancelled {
        fmt.Println("Error: Cancelled")
        return
    }
    
    if currentDate > request.Date {
        fmt.Println("Error: Past reservation")
        return
    }
    
    //restaurant := Restaurants[request.RestaurantId]
 
    //if (currentDate < restaurant.A) || ( currentDate > restaurant.B){
    //    fmt.Println(restaurant.Contact)
		//		return
    //}
		request.IsCancelled = true
		RestaurantRequestDetails[requestId] = request
    fmt.Printf("to:%v %v has been cancelled\n",request.RestaurantId,requestId)
}

func Remove(query string) {
    restaurantId := getStringDetails(query,1)
    _,isPresent:= Restaurants[restaurantId] 
    
    if !isPresent {
        fmt.Println("Error: No such restaurant")
    } else {
        delete(Restaurants, restaurantId)
    }   
}

func Confirm(query string) {
    reservationId:= getStringDetails(query,2)
    
    request:= RestaurantRequestDetails[reservationId]
    
    if request.RequestId != reservationId {
        fmt.Println("Error: No such reservation ID")
				return 
		}
    
    if request.IsRejected {
        fmt.Println("Error: Already rejected")
				return
		}
    
    if request.IsConfirmed {
        fmt.Println("Error: Already confirmed")
				return
		}
    
    if request.IsCancelled {
      fmt.Println("Error: Already cancelled")
			return
		}
    request.IsConfirmed = true
		RestaurantRequestDetails[reservationId] = request
		
    fmt.Printf("to:%v %v has been confirmed\n",request.UserId, reservationId)
}

func RestaurantExists(restaurantId string) bool {
    _,isPresent := Restaurants[restaurantId]
    return isPresent
}

func IsInsideOfReservationPeriod(restaurantId string, requestId string, day int) bool {
    counter:= 0
		if day > 7 {
			day = day - 7
		}

		//  if isNextDay {
		//  	day = day	
		// } else {
		// 	day = day - 1
		// }

    key:=  restaurantId+"_"+strconv.Itoa(day-1)
		_,isPresent := RestaurantTimingDetails[key]
		//fmt.Printf("RestaurantTimingDetails[key] is %v\n", RestaurantTimingDetails[key].Time)
		//fmt.Printf("Key %s is present  \n", key)
    if !isPresent {
			return false
		}


		//fmt.Printf("Inside IsInsideOfReservationPeriod() method where key is %s and day is %d",key, day)
		restaurantTimingDetails:= RestaurantTimingDetails[key].Time
    restaurantRequestDetails:= RestaurantRequestDetails[requestId].Time
    var indexPosition int = -1
        indexPosition = strings.Index(restaurantTimingDetails, " ")    
    
    if indexPosition > -1 {
        times:= strings.Split(restaurantTimingDetails, " ")
        
        for i:=0;i<len(times);i++ {
    hours,_:= strconv.Atoi(strings.Split(restaurantRequestDetails, ":")[0])
    tempTime := strings.Split(times[i], "-")[0]
    minHours,_ := strconv.Atoi(strings.Split(tempTime, ":")[0])
    maxHours,_ := strconv.Atoi(strings.Split(tempTime, ":")[1]) 
            if (hours >= minHours) && (hours <= maxHours) {
                counter++;
            }    
        }
        
        if counter >0 {//} == len(times) {
            return true
        }
    } else {
    hours,_:= strconv.Atoi(strings.Split(restaurantRequestDetails, ":")[0])
    tempTime := strings.Split(restaurantTimingDetails, "-")
    minHours,_ := strconv.Atoi(strings.Split(tempTime[0], ":")[0])
    maxHours,_ := strconv.Atoi(strings.Split(tempTime[1], ":")[0])

        if (hours >= minHours) && (hours <= maxHours) {
            return true
        }
    }
    
    return false
}

func IsRestaurantClosed(restaurantId string) bool{
    //index:= day - 1
    return len(RestaurantTimings[restaurantId].Time) > 0 //[index] == string("-")
}

func IsNotTooFewOrTooManyPeople(restaurantId string, requestId string) bool {
    restaurantNoOfPeople:= Restaurants[restaurantId]
    requestNoOfPeople:= RestaurantRequestDetails[requestId].NoOfPeople
    
    if (requestNoOfPeople >= restaurantNoOfPeople.MinimumReservation) || (requestNoOfPeople <= restaurantNoOfPeople.MaximumReservation){
        return true
    }
    
    return false
}

func IsReservationsTemporarilyClosed(restaurantId string) bool{
    return Restaurants[restaurantId].IsStopped
}

func Request(query string) {
		//fmt.Println("Inside Request() method\n")
    reservationId:= getStringDetails(query,1)
    userId:= getStringDetails(query,2)
    restaurantId:= getStringDetails(query,3)
    date,_ := strconv.Atoi(getStringDetails(query,4))
		//fmt.Printf("In Request() method date is %v\n",date)
    time:=getStringDetails(query,5)
    noOfPeople,_ := strconv.Atoi(getStringDetails(query,6))
    //fmt.Println("Inside Request() method after assigning params",restaurantId)
    
		if !RestaurantExists(restaurantId) {
        fmt.Println("Error: No such restaurant")
        return
    }
    //fmt.Println("Inside Request() method after !RestaurantExists(restaurantId)")
    
    if !IsInsideOfReservationPeriod(restaurantId, reservationId, date) {
        fmt.Println("Error: Outside of reservation period")
        return 
    }
    //fmt.Println("Inside Request() method after !IsInsideOfReservationPeriod(restaurantId, reservationId, date)")

    if IsRestaurantClosed(restaurantId){
        fmt.Println("Error: Closed")
        return
    }
    
    if !IsNotTooFewOrTooManyPeople(restaurantId, reservationId) {
        fmt.Println("Error: Too many or too few people")
        return
    }
    
    if IsReservationsTemporarilyClosed(restaurantId) {
        fmt.Println("Error: Reservations temporarily closed")
        return 
    }
    
    fmt.Printf("to:%v Received a reservation request: %v %v %v %v %v\n",restaurantId,reservationId,userId,date,time,noOfPeople)
}

func getStringDetails(text string, index int) string {
    return strings.Fields(text)[index]
}

func AddRestaurantDetails(text string) {
    restaurantId:= getStringDetails(text,0)
    contact:= getStringDetails(text,1)
    a,_:= strconv.Atoi(getStringDetails(text,2))
    b,_:= strconv.Atoi(getStringDetails(text,3))
    min,_:= strconv.Atoi(getStringDetails(text,4))
    max,_:= strconv.Atoi(getStringDetails(text,5))

    Restaurants[restaurantId] = Restaurant{ RestaurantId: restaurantId, Contact: contact, A: a, B: b, MinimumReservation: min, MaximumReservation: max,IsStopped : false  }
}

func AddRestaurantTimingDetails(restaurantId string, index int, text string) {
    //fmt.Println("Inside AddRestaurantTimingDetails()")
		if index < 0 {
			index = index * -1
		}

		//time:= getStringDetails(text, )

    tempIndex:= restaurantId+"_"+strconv.Itoa(index)
		//fmt.Println("Index is ", tempIndex) to comment
		// if restaurantId == "ssqteaj"{
		// 	fmt.Printf("Time to be added to index : %s is value : %s\n",restaurantId, text)
		// }

    RestaurantTimingDetails[tempIndex] = RestaurantTimingDetail{ RestaurantId : restaurantId, Time : text }
    
    //fmt.Println(RestaurantTimingDetails[tempIndex]);
}

func AddRestaurantRequestDetails(restaurantId string, userId string, reservationId string, date int, time string, noOfPeople int){
    _,isPresent:= RestaurantRequestDetails[reservationId]
    
    if(!isPresent) {
        RestaurantRequestDetails[reservationId] = RestaurantRequestDetail{RequestId :reservationId, UserId: userId, RestaurantId: restaurantId, Date: date, Time: time, NoOfPeople: noOfPeople, IsPending: true, IsConfirmed: false, IsRejected: false }     
    }
}

func executeQuery(query string) {
	command:= getStringDetails(query,0)
	switch command {
        case "CONFIRM":
        Confirm(query);
        
        case "REQUEST":
				restaurantId:= getStringDetails(query, 3) 
				requestId:= getStringDetails(query, 1)
				userId:= getStringDetails(query, 2)
				date,_:= strconv.Atoi(getStringDetails(query, 4))
				time:= getStringDetails(query, 5)
				noOfPeople,_:= strconv.Atoi(getStringDetails(query, 6))
				AddRestaurantRequestDetails(restaurantId, userId, requestId, date, time, noOfPeople)
        Request(query);
       
        case "NEXT_DAY":
        NextDay();
        
        case "LIST":
        List();
        
        case "STOP":
        Stop(query);
        
        case "REJECT":
        Reject(query);
        
        case "CANCEL":
        Cancel(query);
        
        case "REMOVE":
        Remove(query);
        
        //default :
        //fmt.Println("Invalid query");
        // to add break 
    }
}

func main() {
	InitCommands()
	restaurantCounter:=0
	noOfRestaurant:=0
	restaurantId:=""
	tempIndex:=0
	lines := getStdin()
	for i, v := range lines {
		key := getStringDetails(v,0)
		_,isPresentCommand:= commands[key]		
		//fmt.Printf("line[%d] %v\n",i, v)
		if i == 0 {
		 	noOfRestaurant,_ =  strconv.Atoi(v)
				
		 	//fmt.Printf("noOfRestaurant is %d", noOfRestaurant)
		} else if isPresentCommand {
			//fmt.Printf("Inside executeQuery query :- %v\n", v)
			executeQuery(v)
		} else	if  i == 1  {
			AddRestaurantDetails(v)
			restaurantId = getStringDetails(v,0)
			//fmt.Printf("Restaurant id is %s\n", restaurantId)
			restaurantCounter++
		} else if i<=8 {
				//for i:=0; i<8;i++ {
					AddRestaurantTimingDetails(restaurantId, i-2, v) 
					//fmt.Printf("Inside add restaurant timing for %d and index is %d\n",i, i-2)	
				//}			
		} else if (listCounter == 0) && (noOfRestaurant > 1) && (restaurantCounter <= noOfRestaurant) {
			//fmt.Printf("inside if (noOfRestaurant > 1) && (restaurantCounter < noOfRestaurant) and i is :%d\n",i)
			if (i-1)%8 == 0 {
				//fmt.Printf("inside %d restaurant add\n",restaurantCounter);
				AddRestaurantDetails(v)	
				restaurantId = getStringDetails(v,0)
				//fmt.Printf("Restaurant id is %s\n", restaurantId)
				//fmt.Printf("Inside add restaurant for %d",restaurantCounter)
				tempIndex = i
				restaurantCounter++
				//fmt.Printf("restaurantCounter is %d\n",restaurantCounter)
			} else if (i <= (i+7)) {//&& (restaurantCounter < noOfRestaurant) {
				//for i:=0; i<8;i++ {
					//fmt.Printf("inside restaurant timing add for %d and restaurantId is %v\n", i, restaurantCounter,restaurantId );
					AddRestaurantTimingDetails(restaurantId, i - tempIndex - 1, v) 	
					//fmt.Printf("Inside add restaurant timing for %d",i)	
				//}		
			}
		} else {
		//fmt.Printf("Inside list section\n")
		if listCounter > 0 {
			if listCounter == 1 {
				tempIndex = i
				//fmt.Println("AddRestaurantDetails for ",i)
				AddRestaurantDetails(v)
				listCounter++
				restaurantId = getStringDetails(v,0)
				//fmt.Printf("Restaurant id is %s\n", restaurantId)
			} else if listCounter <= 9 {
				//fmt.Println("AddRestaurantTimingDetails for ",i, listCounter - 2) 	
				AddRestaurantTimingDetails(restaurantId, listCounter - 2, v)
				listCounter++
			} 

			if listCounter == 9 {
				listCounter = 0
			}
		} 
		}
	}
	
	//fmt.Printf("===========================\n")
	//fmt.Printf(" %s\n", lines[2])


}

func getStdin() []string {
	stdin, _ := ioutil.ReadAll(os.Stdin)
	return strings.Split(strings.TrimRight(string(stdin), "\n"), "\n")
}
