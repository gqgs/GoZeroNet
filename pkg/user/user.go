package user

import (
	"encoding/json"
	"os"
	"path"

	"github.com/gqgs/go-zeronet/pkg/config"
	"github.com/gqgs/go-zeronet/pkg/lib/crypto"
)

type User interface {
	AuthAddress(addr string) string
	CertUserID(addr string) string
	SiteSettings(addr string) map[string]interface{}
	SetSiteSettings(addr string, settings map[string]interface{}) error
	GlobalSettings() GlobalSettings
}

type manager struct {
	users map[string]*user
}

type Manager interface {
	User() User
}

func (m *manager) User() User {
	// For now return the first user
	for _, user := range m.users {
		return user
	}

	// user doesn't exist
	user := new(user)
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

func loadUserSettingsFromFile() (map[string]*user, error) {
	userFilePath := path.Join(config.DataDir, "users.json")
	file, err := os.Open(userFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]*user), nil
		}
		return nil, err
	}
	defer file.Close()

	users := make(map[string]*user)
	return users, json.NewDecoder(file).Decode(&users)
}

type Cert struct {
	AuthAddress    string `json:"auth_address"`
	AuthPrivatekey string `json:"auth_privatekey"`
	AuthType       string `json:"auth_type"`
	AuthuserName   string `json:"auth_user_name"`
	CertSign       string `json:"cert_sign"`
}

type Site struct {
	AuthAddress    string                 `json:"auth_address"`
	AuthPrivatekey string                 `json:"auth_privatekey"`
	Settings       map[string]interface{} `json:"settings"`
}

type GlobalSettings struct {
	Theme          string `json:"theme"`
	UseSystemTheme bool   `json:"use_system_theme"`
}

type user struct {
	addr       string
	Certs      map[string]Cert  `json:"certs"`
	MasterSeed string           `json:"master_seed"`
	Settings   GlobalSettings   `json:"settings"`
	Sites      map[string]*Site `json:"sites"`
}

func (u *user) AuthAddress(addr string) string {
	if u.Sites[addr] == nil {
		return ""
	}
	return u.Sites[addr].AuthAddress
}

func (u *user) CertUserID(addr string) string {
	// TODO: implement me
	return ""
}

func (u *user) SiteSettings(addr string) map[string]interface{} {
	if u.Sites[addr] == nil {
		return nil
	}
	return u.Sites[addr].Settings
}

func (u *user) GlobalSettings() GlobalSettings {
	return u.Settings
}

func (u *user) SetSiteSettings(addr string, settings map[string]interface{}) error {
	users, err := loadUserSettingsFromFile()
	if err != nil {
		return err
	}

	if u.Sites[addr] == nil {
		u.Sites[addr] = new(Site)
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
