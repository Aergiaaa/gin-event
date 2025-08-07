package main

import (
	"net/http"
	"strconv"

	"github.com/Aergiaaa/gin-event/internal/database"

	"github.com/gin-gonic/gin"
)

// GetAttendeesForEvent returns all attendees for a given event
//
//	@Summary			Returns all attendees for a given event
//	@Description	Returns all attendees for a given event
//	@Tags				attendees
//	@Accept			json
//	@Produce			json
//	@Param			id		path		int	true	"Event ID"
//	@Success			200	{object}	[]database.User
//	@Router			/api/v1/events/{id}/attendees [get]
func (app *app) getAttendeesForEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event Id"})
	}

	users, err := app.models.Attendees.GetByEvent(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving attendees"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// AddAttendeeToEvent adds an attendee to an event
// @Summary			Adds an attendee to an event
// @Description	Adds an attendee to an event
// @Tags				attendees
// @Accept			json
// @Produce			json
// @Param			id			path		int	true	"Event ID"
// @Param			userId	path		int	true	"User ID"
// @Success			201		{object}	database.Attendee
// @Router			/api/v1/events/{id}/attendees/{userId} [post]
// @Security		BearerAuth
func (app *app) addAttendeeToEvent(c *gin.Context) {
	eventId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event Id"})
		return
	}

	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user Id"})
		return
	}

	event, err := app.models.Events.Get(eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving event"})
		return
	}
	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	user := app.getUserFromContext(c)
	if event.OwnerID != user.Id {
		c.JSON(http.StatusForbidden,
			gin.H{"error": "You are not allowed to add attendees to this event"})
		return
	}

	userToAdd, err := app.models.Users.Get(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving user"})
		return
	}
	if userToAdd == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	existingAttendee, err := app.models.Attendees.GetByEventAndUser(event.Id, userToAdd.Id)
	if existingAttendee != nil {
		c.JSON(http.StatusConflict,
			gin.H{"error": "User is already an attendee of this event"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": "Error retrieving attendee"})
		return
	}

	attendee := &database.Attendee{
		EventId: event.Id,
		UserId:  userToAdd.Id,
	}

	_, err = app.models.Attendees.Insert(attendee)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": "Error adding attendee to event"})
		return
	}

	c.JSON(http.StatusCreated, attendee)
}

// DeleteAttendeeFromEvent deletes an attendee from an event
// @Summary			Deletes an attendee from an event
// @Description	Deletes an attendee from an event
// @Tags				attendees
// @Accept			json
// @Produce			json
// @Param			id			path		int	true	"Event ID"
// @Param			userId	path		int	true	"User ID"
// @Success			204
// @Router			/api/v1/events/{id}/attendees/{userId} [delete]
// @Security		BearerAuth
func (app *app) deleteAttendeeFromEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event Id"})
		return
	}

	event, err := app.models.Events.Get(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving event"})
		return
	}
	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user Id"})
		return
	}

	err = app.models.Attendees.Delete(userId, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete attendee"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// GetEventsByAttendee returns all events for a given attendee
//
//	@Summary			Returns all events for a given attendee
//	@Description	Returns all events for a given attendee
//	@Tags				attendees
//	@Accept			json
//	@Produce			json
//	@Param			id		path		int	true	"Attendee ID"
//	@Success			200	{object}	[]database.Event
//	@Router			/api/v1/attendees/{id}/events [get]
func (app *app) getEventsByAttendee(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid attendees Id"})
		return
	}

	events, err := app.models.Attendees.GetEventsByUserId(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving events"})
		return
	}

	c.JSON(http.StatusOK, events)
}
