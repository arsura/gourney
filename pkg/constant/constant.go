package constant

import "time"

type CorrelationIdType string

const (
	REQUEST_ID_KEY CorrelationIdType = "request_id"
)

type RoutingKeyPrefix string

const (
	ROUTING_KEY_PREFIX_LOGS RoutingKeyPrefix = "logs"
)

type TargetEvent string

const (
	TARGET_POST TargetEvent = "post"
)

type EventAction string

const (
	ACTION_CREATE EventAction = "create"
	ACTION_UPDATE EventAction = "update"
	ACTION_DELETE EventAction = "delete"
)

const ROUTING_KEY_LOGS = "logs.*.*" // e.g. logs.post.create

const MAX_TEMP_MESSAGE_SIZE = 5
const TICK_TIME_TO_WRITE_LOGS = 1 * 60 * time.Second
