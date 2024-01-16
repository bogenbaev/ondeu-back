package modules

type AppConfigs struct {
	Port          string
	LogLevel      string
	Keycloak      *Keycloak
	Database      *Postgre
	ObjectStorage *ObjectStorage
}

type ObjectStorage struct {
	Endpoint     string
	Bucket       string
	ClientName   string
	ClientSecret string
	ClientKey    string
}

type Postgre struct {
	Host     string
	Port     int
	Username string
	Password string
	DBName   string
	SSLMode  string
}

type Keycloak struct {
	Host              string
	Realm             string
	ClientID          string
	AdminClientID     string
	AdminClientSecret string
}
