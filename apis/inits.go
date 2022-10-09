package apis

import "happy_day/application"

func initCreateReservationHandler() application.CreateReservationHandler {
	return application.CreateReservationHandler{}
}

func initChangeReservationHandler() application.ChangeReservationHandler {
	return application.ChangeReservationHandler{}
}

func initDeleteReservationHandler() application.DeleteReservationHandler {
	return application.DeleteReservationHandler{}
}

func initGetReservationByIdHandler() application.GetReservationByIdHandler {
	return application.GetReservationByIdHandler{}
}

func initGetAllReservationHandler() application.GetAllReservationHandler {
	return application.GetAllReservationHandler{}
}
