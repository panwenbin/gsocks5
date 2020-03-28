package consts

const VERSION4 = uint8(4)
const VERSION5 = uint8(5)

const (
	METHOD_NO_AUTHENTICATION_REQUIRED = 0x00
	METHOD_GSSAPI                     = 0x01
	METHOD_USERNAME_PASSWORD          = 0x02
	METHOD_NO_ACCEPTABLE_METHODS      = 0xFF
	// 0x03 to 0x7F IANA ASSIGNED
	// 0x80 to 0xFE RESERVED FOR PRIVATE METHODS
)

const (
	CMD_CONNECT       = 0x01
	CMD_BIND          = 0x02
	CMD_UDP_ASSOCIATE = 0x03
)

const RSV = 0x00

const (
	ATYP_IPv4   = 0x01
	ATYP_DOMAIN = 0x03
	ATYP_IPv6   = 0x04
)

const (
	REP_SUCCEEDED                         = 0x00
	REP_GENERAL_SOCKS_SERVER_FAILURE      = 0x01
	REP_CONNECTION_NOT_ALLOWED_BY_RULESET = 0x02
	REP_NETWORK_UNREACHABLE               = 0x03
	REP_HOST_UNREACHABLE                  = 0x04
	REP_CONNECTION_REFUSED                = 0x05
	REP_TTL_EXPIRED                       = 0x06
	REP_COMMAND_NOT_SUPPORTED             = 0x07
	REP_ADDRESS_TYPE_NOT_SUPPORTED        = 0x08
	// 0x09 to 0xFF unassigned
)
