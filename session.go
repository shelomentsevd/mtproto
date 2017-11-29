package mtproto

import (
	"errors"
	"os"
	"math/rand"
	"time"
)

// Session storage interface
type ISession interface {
	// Load is deserialization method
	Load() error
	// Save is serialization method
	Save() error

	IsIPv6() bool
	// IsEncrypted returns true if AuthKey, ServerSalt and SessionID fields aren't empty
	IsEncrypted() bool

	GetAddress() string
	GetAuthKey() []byte
	GetAuthKeyHash() []byte
	GetServerSalt() []byte
	GetSessionID() int64

	SetAddress(string)
	SetAuthKey([]byte)
	SetAuthKeyHash([]byte)
	SetServerSalt([]byte)
	SetSessionID(int64)

	UseIPv6(bool)
	Encrypted(bool)
}

type Session struct {
	// TODO: ReaderWriter interface
	file *os.File

	address     string
	authKey     []byte
	authKeyHash []byte
	serverSalt  []byte
	sessionId   int64
	useIPv6     bool
	encrypted   bool
}

func NewSession(file *os.File) ISession {
	session := &Session{
		file: file,
	}

	rand.Seed(time.Now().UnixNano())
	session.SetSessionID(rand.Int63())

	return session
}

func (s *Session) Load() error {
	// TODO: Magic number
	buffer := make([]byte, 1024*4)
	n, _ := s.file.ReadAt(buffer, 0)
	if n <= 0 {
		return errors.New("New session")
	}

	decoder := NewDecodeBuf(buffer)
	s.authKey = decoder.StringBytes()
	s.authKeyHash = decoder.StringBytes()
	s.serverSalt = decoder.StringBytes()
	s.address = decoder.String()
	s.useIPv6 = false
	if decoder.UInt() == 1 {
		s.useIPv6 = true
	}

	if decoder.err != nil {
		return decoder.err
	}

	return nil
}

func (s Session) Save() error {
	// TODO: Magic number
	buffer := NewEncodeBuf(1024)
	buffer.StringBytes(s.authKey)
	buffer.StringBytes(s.authKeyHash)
	buffer.StringBytes(s.serverSalt)
	buffer.String(s.address)

	var useIPv6UInt uint32
	if s.useIPv6 {
		useIPv6UInt = 1
	}
	buffer.UInt(useIPv6UInt)

	err := s.file.Truncate(0)
	if err != nil {
		return err
	}

	_, err = s.file.WriteAt(buffer.buf, 0)
	if err != nil {
		return err
	}

	return nil
}

func (s Session) IsIPv6() bool {
	return s.useIPv6
}

func (s Session) IsEncrypted() bool {
	return s.encrypted
}

func (s Session) GetAddress() string {
	return s.address
}

func (s Session) GetAuthKey() []byte {
	return s.authKey
}

func (s Session) GetAuthKeyHash() []byte {
	return s.authKeyHash
}

func (s Session) GetServerSalt() []byte {
	return s.serverSalt
}

func (s Session) GetSessionID() int64 {
	return s.sessionId
}

func (s *Session) SetAddress(address string) {
	s.address = address
}

func (s *Session) SetAuthKey(authKey []byte) {
	s.authKey = make([]byte, len(authKey))
	copy(s.authKey, authKey)
}

func (s *Session) SetAuthKeyHash(authKeyHash []byte) {
	s.authKeyHash = make([]byte, len(authKeyHash))
	copy(s.authKeyHash, authKeyHash)
}

func (s *Session) SetServerSalt(salt []byte) {
	s.serverSalt = make([]byte, len(salt))
	copy(s.serverSalt, salt)
}

func (s *Session) SetSessionID(ID int64) {
	s.sessionId = ID
}

func (s *Session) UseIPv6(useIPv6 bool) {
	s.useIPv6 = useIPv6
}

func (s *Session) Encrypted(encrypted bool) {
	s.encrypted = encrypted
}
