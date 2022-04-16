package libomni

type Omni interface {
	Admiter
	Codec
	Autorizer
}

type Session struct {
	ID string
}

type Admiter interface {
	Admit(Session, []byte) bool
}

type Codec interface {
	Encode(Session, EncodeInput) ([]byte, error)
	Decode(Session, []byte) (*Signal, error)
}

type Encoder interface {
	Encode(Session, EncodeInput) ([]byte, error)
}

type Decoder interface {
	Decode(Session, []byte) (*Signal, error)
}

type Autorizer interface {
	Authorize(Session, Device, map[string]string) (bool, error)
}

type Closer interface {
	Close() error
}

type Device struct {
	ID   string
	IMEI string
}

type Signal struct {
	Device               *Device
	AuthorizationRequest *AuthorizationRequest
	PositionUpdate       *PositionUpdate
	Close                *bool
}

type AuthorizationRequest struct {
	Credentials map[string]string
}

type PositionUpdate struct {
}

type EncodeInput struct {
	PositionUpdateResponse *PositionUpdateResponse
	AuthorizationResponse  *AuthorizationResponse
	Ignite                 *Ignite
}

type PositionUpdateResponse struct {
}

type AuthorizationResponse struct {
}

type Ignite struct {
}

type MissingProviderError struct {
	err error
}

func (e MissingProviderError) Error() string {
	return e.err.Error()
}

type ParseError struct {
	err error
}

func (e ParseError) Error() string {
	return e.err.Error()
}
