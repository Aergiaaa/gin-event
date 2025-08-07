package main

import (
	"net/http"
	"strconv"

	"github.com/Aergiaaa/gin-event/internal/database"

	"github.com/gin-gonic/gin"
)

// GetEvent returns a single event
//
//	@Summary			Returns a single event
//	@Description	Returns a single event
//	@Tags				events
//	@Accept			json
//	@Produce			json
//	@Param			id		path		int	true	"Event ID"
//	@Success			200	{object}	database.Event
//	@Router			/api/v1/events/{id} [get]
func (app *app) getEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	event, err := app.models.Events.Get(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return
	}
	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	c.JSON(http.StatusOK, event)
}

// GetEvents returns all events
//
//	@Summary			Returns all events
//	@Description	Returns all events
//	@Tags				events
//	@Accept			json
//	@Produce			json
//	@Success			200		{object}		[]database.Event
//	@Router			/api/v1/events [get]
func (app *app) getAllEvents(c *gin.Context) {
	events, err := app.models.Events.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve events"})
		return
	}

	c.JSON(http.StatusOK, events)
}

// CreateEvent creates a new event
//
//	@Summary			Creates a new event
//	@Description	Creates a new event
//	@Tags				events
//	@Accept			json
//	@Produce			json
//	@Param			event	body		database.Event	true	"Event"
//	@Success			201	{object}	database.Event
//	@Router			/api/v1/events [post]
//	@Security		BearerAuth
func (app *app) createEvent(c *gin.Context) {
	var event database.Event

	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := app.getUserFromContext(c)
	event.OwnerID = user.Id

	if err := app.models.Events.Insert(&event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event"})
		return
	}

	c.JSON(http.StatusCreated, event)
}

// UpdateEvent updates an existing event
//
//	@Summary			Updates an existing event
//	@Description	Updates an existing event
//	@Tags				events
//	@Accept			json
//	@Produce			json
//	@Param			id		path		int				true	"Event ID"
//	@Param			event	body		database.Event	true	"Event"
//	@Success			200	{object}	database.Event
//	@Router			/api/v1/events/{id} [put]
//	@Security		BearerAuth
func (app *app) updateEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	existingEvent, err := app.models.Events.Get(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return
	}
	if existingEvent == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	user := app.getUserFromContext(c)
	if existingEvent.OwnerID != user.Id {
		c.JSON(http.StatusForbidden,
			gin.H{"error": "You do not have permission to update this event"})
		return
	}

	updatedEvent := &database.Event{}
	if err := c.ShouldBindJSON(updatedEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedEvent.Id = id
	if err := app.models.Events.Update(updatedEvent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update event"})
		return
	}

	c.JSON(http.StatusOK, updatedEvent)
}

// DeleteEvent deletes an existing event
//
//	@Summary			Deletes an existing event
//	@Description	Deletes an existing event
//	@Tags				events
//	@Accept			json
//	@Produce			json
//	@Param			id		path	int	true	"Event ID"
//	@Success			204
//	@Router			/api/v1/events/{id} [delete]
//	@Security		BearerAuth
func (app *app) deleteEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
	}

	user := app.getUserFromContext(c)
	event, err := app.models.Events.Get(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return
	}
	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	if event.OwnerID != user.Id {
		c.JSON(http.StatusForbidden,
			gin.H{"error": "You do not have permission to delete this event"})
		return
	}

	if err := app.models.Events.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete event"})
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "Event deleted successfully"})
}
