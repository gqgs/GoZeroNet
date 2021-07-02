package user

import (
	"encoding/json"
	"errors"
	"os"
	"path"

	"github.com/gqgs/go-zeronet/pkg/config"
	"github.com/gqgs/go-zeronet/pkg/lib/crypto"
)

var (
	ErrCertNotChanged    = errors.New("user: cert not changed")
	ErrCertAlreadyExists = errors.New("user: cert already exists")
)

type manager struct {
	users map[string]*User
}

type Manager interface {
	User() *User
}

func (m *manager) User() *User {
	// For now return the first user
	for _, user := range m.users {
		return user
	}

	// user doesn't exist
	user := new(User)
	user.Sites = make(map[string]*Site)
	user.MasterSeed = crypto.NewPrivateKey(crypto.Hex)
	addr, _ := crypto.PrivateKeyToAddress(user.MasterSeed)
	user.addr = addr
	m.users[addr] = user
	return user
}

func NewManager() (*manager, error) {
	users, err := loadUserSettingsFromFile()
	if err != nil {
		return nil, err
	}
	for addr, user := range users {
		user.addr = addr
	}

	return &manager{users}, nil
}

func loadUserSettingsFromFile() (map[string]*User, error) {
	userFilePath := path.Join(config.DataDir, "users.json")
	file, err := os.Open(userFilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return make(map[string]*User), nil
		}
		return nil, err
	}
	defer file.Close()

	users := make(map[string]*User)
	return users, json.NewDecoder(file).Decode(&users)
}

type Cert struct {
	AuthAddress    string `json:"auth_address"`
	AuthPrivatekey string `json:"auth_privatekey"`
	AuthType       string `json:"auth_type"`
	AuthUserName   string `json:"auth_user_name"`
	CertSign       string `json:"cert_sign"`
}

type Site struct {
	AuthAddress    string                 `json:"auth_address"`
	AuthPrivatekey string                 `json:"auth_privatekey"`
	Cert           string                 `json:"cert"`
	Settings       map[string]interface{} `json:"settings"`
}

type GlobalSettings struct {
	Theme          string `json:"theme"`
	UseSystemTheme bool   `json:"use_system_theme"`
}

type User struct {
	addr       string
	Certs      map[string]Cert  `json:"certs"`
	MasterSeed string           `json:"master_seed"`
	Settings   GlobalSettings   `json:"settings"`
	Sites      map[string]*Site `json:"sites"`
}

func (u *User) AuthAddress(addr string) string {
	if u.Sites[addr] == nil {
		return ""
	}
	return u.Sites[addr].AuthAddress
}

func (u *User) CertUserID(addr string) string {
	if u.Sites[addr] == nil {
		return ""
	}

	if u.Certs[u.Sites[addr].Cert].AuthUserName == "" {
		return ""
	}

	return u.Certs[u.Sites[addr].Cert].AuthUserName + "@" + u.Sites[addr].Cert
}

func (u *User) SiteSettings(addr string) map[string]interface{} {
	if u.Sites[addr] == nil {
		return nil
	}
	return u.Sites[addr].Settings
}

func (u *User) GlobalSettings() GlobalSettings {
	return u.Settings
}

func (u *User) SetSiteSettings(addr string, settings map[string]interface{}) error {
	users, err := loadUserSettingsFromFile()
	if err != nil {
		return err
	}

	if u.Sites[addr] == nil {
		u.Sites[addr] = new(Site)
		if u.Sites[addr].AuthPrivatekey, err = crypto.AuthPrivateKey(u.MasterSeed, addr); err != nil {
			return err
		}
		if u.Sites[addr].AuthAddress, err = crypto.PrivateKeyToAddress(u.Sites[addr].AuthPrivatekey); err != nil {
			return err
		}
	}

	u.Sites[addr].Settings = settings
	users[u.addr] = u

	data, err := json.Marshal(users)
	if err != nil {
		return err
	}

	userFilePath := path.Join(config.DataDir, "users.json")
	return os.WriteFile(userFilePath, data, os.ModePerm)
}

func (u *User) SaveSettings() error {
	users, err := loadUserSettingsFromFile()
	if err != nil {
		return err
	}
	users[u.addr] = u

	data, err := json.Marshal(users)
	if err != nil {
		return err
	}

	userFilePath := path.Join(config.DataDir, "users.json")
	return os.WriteFile(userFilePath, data, os.ModePerm)
}

func (u *User) UpdateCert(addr, cert string) error {
	if u.Sites[addr] == nil {
		return errors.New("site not found")
	}
	u.Sites[addr].Cert = cert
	return u.SaveSettings()
}

func (u *User) AddCert(addr, domain, authType, authUserName, certSign string, force bool) error {
	if u.Sites[addr] == nil {
		return errors.New("site not found")
	}
	cert, exists := u.Certs[domain]
	if exists && !force {
		if cert.AuthType == authType &&
			cert.AuthUserName == authUserName &&
			cert.CertSign == certSign {
			return ErrCertNotChanged
		}
		return ErrCertAlreadyExists
	}

	cert.AuthPrivatekey = u.Sites[addr].AuthPrivatekey
	cert.AuthAddress = u.Sites[addr].AuthAddress
	cert.AuthType = authType
	cert.AuthUserName = authUserName
	cert.CertSign = certSign

	u.Certs[domain] = cert
	u.Sites[addr].Cert = domain

	return u.SaveSettings()
}
