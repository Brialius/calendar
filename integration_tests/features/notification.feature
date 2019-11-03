Feature: Test notification service

	Scenario: Notification for Created Event
		Given there is user "test_user"
		And there is server "calendar-service:8080"
		And there is MQ server "amqp://queue_user:queue-super-password@rabbit:5672/"
		And MQ queue "test_notification"
		And MQ route key "notification.tasks"
		And MQ exchange "calendar"
		When I create event
		"""
		{
			"title":"test_event_create_get",
			"text":"Test event: Notification check",
			"startTime":"2019-10-01T00:00:00Z",
			"endTime":"2019-10-02T00:00:00Z"
		}
		"""
		And I get task from queue
		Then Event should be the same as created
		And I delete event by previous id
