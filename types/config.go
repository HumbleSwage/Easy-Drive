package types

type Server struct {
	Domain   string `yaml:"domain"`
	Version  string `yaml:"version"`
	AppEnv   string `yaml:"appEnv"`
	HttpPort string `yaml:"httpPort"`
}

type Mysql struct {
	Dialect  string `yaml:"dialect"`
	DbHost   string `yaml:"dbHost"`
	DbPort   string `yaml:"dbPort"`
	DbName   string `yaml:"dbName"`
	UserName string `yaml:"userName"`
	Password string `yaml:"password"`
	Charset  string `yaml:"charset"`
}

type Redis struct {
	RedisDbName   int    `yaml:"redisDbName"`
	RedisHost     string `yaml:"redisHost"`
	RedisPort     string `yaml:"redisPort"`
	RedisPassword string `yaml:"redisPassword"`
	RedisNetwork  string `yaml:"redisNetwork"`
}

type Email struct {
	ValidEmail   string `yaml:"validEmail"`
	SmtpHost     string `yaml:"smtpHost"`
	SmtpEmail    string `yaml:"smtpEmail"`
	SmtpPassword string `yaml:"smtpPassword"`
}

type PhotoPath struct {
	Host string `yaml:"host"`
}

type UploadPath struct {
	Host        string `yaml:"host"`
	TempPath    string `yaml:"tempPath"`
	AvatarPath  string `yaml:"avatarPath"`
	VideoPath   string `yaml:"videoPath"`
	MusicPath   string `yaml:"musicPath"`
	ImagePath   string `yaml:"imagePath"`
	DocPath     string `yaml:"docPath"`
	ProgramPath string `yaml:"programPath"`
	ZipPath     string `yaml:"zipPath"`
	OthersPath  string `yaml:"othersPath"`
}
