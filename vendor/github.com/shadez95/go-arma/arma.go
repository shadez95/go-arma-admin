package arma

import (
	"log"
	"os"

	"os/exec"
	"strings"
)

var configs = "configs"

// Platform type represents what platform the server will run on
type Platform string

const (
	// Linux is of type Platform and should be used when running in pure Linux
	Linux Platform = "linux"
	// Windows is of type Platform and used when running in Windows
	Windows Platform = "windows"
	// Wine is of type Platform and should be used when running in linux but using wnie
	Wine Platform = "wine"
)

// BaseConfig is base struct for server and headless clients
type BaseConfig struct {
	// Executable is the actual binary file. Specify if you want to use 32 bit binary. Default uses 64 bit binary
	Executable string
	// Path is the path to the ARMA directory. Not the path to the executable
	Path string
	// Platform is either linux, windows, or wine
	Platform Platform
	// Port the server will run on
	Port     string
	NoLogs   bool
	EnableHT bool
	Profiles string
	Mods     []string
	BEPath   string
}

// Server is a struct for an ARMA server
type Server struct {
	AutoInit            bool
	BasicConfig         string
	LoadMissionToMemory bool
	ProfileName         string
	ServerConfig        string
	ServerMod           string
	*BaseConfig
}

// HeadlessClient is a struct that will represent a headless client
type HeadlessClient struct {
	Connect     string
	Password    string
	ProfileName string
	*BaseConfig
}

// NewServer starts a new server instance with defaults set. Only Path and Platform are required in BaseConfig struct
func NewServer(server *Server) *Server {

	// Set defaults for BaseConfig portion of the struct
	server.Port = "2302"
	server.EnableHT = true
	server.Profiles = "profiles"

	// Set defaults for Server portion of the struct
	server.BasicConfig = "Default_basic.cfg"
	server.LoadMissionToMemory = true
	server.ServerConfig = "Default_server.cfg"

	return server
}

// Start an Arma server. Returns path to executable and args used
func (s *Server) Start() chan *os.Process {

	var mods string
	var args []string
	var armaExecutable string

	switch s.Platform {
	case Windows:
		if len(s.Executable) > 0 {
			armaExecutable = strings.Join([]string{s.Path, s.Executable}, "\\")
		} else {
			armaExecutable = strings.Join([]string{s.Path, "arma3server_x64.exe"}, "\\")
		}
		mods = strings.Join(s.Mods, ";")
	case Wine:
		if len(s.Executable) > 0 {
			armaExecutable = strings.Join([]string{s.Path, s.Executable}, "\\")
		} else {
			armaExecutable = strings.Join([]string{s.Path, "arma3server_x64.exe"}, "\\")
		}
		mods = strings.Join(s.Mods, ";")
	case Linux:
		if len(s.Executable) > 0 {
			armaExecutable = strings.Join([]string{s.Path, s.Executable}, "/")
		} else {
			armaExecutable = strings.Join([]string{s.Path, "arma3server_x64"}, "/")
		}
		mods = strings.Join(s.Mods, `\;`)
	default:
		// s.Logger.Error("Platform not specified, should be windows, wine, or linux")
		log.Fatal("Platform not specified, should be windows, wine, or linux")
	}

	args = append(args, "-port="+s.Port)

	args = append(args, "-noSound")

	if s.NoLogs {
		args = append(args, "-noLogs")
	}

	if s.EnableHT {
		args = append(args, "-enableHT")
	}

	if len(s.Profiles) > 0 {
		args = append(args, "-profiles="+s.Profiles)
	}

	if len(s.Mods) > 0 {
		args = append(args, "-mod="+mods)
	}

	if s.AutoInit {
		args = append(args, "-autoInit")
	}

	if len(s.BasicConfig) > 0 {
		args = append(args, "-cfg="+s.BasicConfig)
	}

	if len(s.BEPath) > 0 {
		args = append(args, "-bepath="+s.BEPath)
	}

	if s.LoadMissionToMemory {
		args = append(args, "-loadMissionToMemory")
	}

	if len(s.ProfileName) > 0 {
		args = append(args, "-name="+s.ProfileName)
	}

	if len(s.ServerConfig) > 0 {
		args = append(args, "-config="+s.ServerConfig)
	}

	if len(s.ServerMod) > 0 {
		args = append(args, "-serverMod="+s.ServerMod)
	}

	c := exec.Command(armaExecutable, args...)
	processCh := make(chan *os.Process)
	go func() {
		processCh <- c.Process
	}()
	c.Start()
	return processCh
}

// NewHeadlessClient starts a new instance of HeadlessClient with defaults set. ProfileName are required
func NewHeadlessClient(hc *HeadlessClient) *HeadlessClient {
	hc.Connect = "127.0.0.1"
	hc.Port = "2302"
	return hc
}

// Start an Arma server. Returns path to executable and args used
func (s *HeadlessClient) Start() chan *os.Process {

	var mods string
	var args []string
	var armaExecutable string

	switch s.Platform {
	case Windows:
		if len(s.Executable) > 0 {
			armaExecutable = strings.Join([]string{s.Path, s.Executable}, "\\")
		} else {
			armaExecutable = strings.Join([]string{s.Path, "arma3server_x64.exe"}, "\\")
		}
		mods = strings.Join(s.Mods, ";")
	case Wine:
		if len(s.Executable) > 0 {
			armaExecutable = strings.Join([]string{s.Path, s.Executable}, "\\")
		} else {
			armaExecutable = strings.Join([]string{s.Path, "arma3server_x64.exe"}, "\\")
		}
		mods = strings.Join(s.Mods, ";")
	case Linux:
		if len(s.Executable) > 0 {
			armaExecutable = strings.Join([]string{s.Path, s.Executable}, "/")
		} else {
			armaExecutable = strings.Join([]string{s.Path, "arma3server_x64"}, "/")
		}
		mods = strings.Join(s.Mods, `\;`)
	default:
		// s.Logger.Error("Platform not specified, should be windows, wine, or linux")
		log.Fatal("Platform not specified, should be windows, wine, or linux")
	}

	args = append(args, "-client")
	args = append(args, "-connect="+s.Connect)
	args = append(args, "-port="+s.Port)
	args = append(args, "-noSound")

	if s.NoLogs {
		args = append(args, "-noLogs")
	}

	if s.EnableHT {
		args = append(args, "-enableHT")
	}

	if len(s.Profiles) > 0 {
		args = append(args, "-profiles="+s.Profiles)
	}

	if len(s.Mods) > 0 {
		args = append(args, "-mod="+mods)
	}

	if len(s.BEPath) > 0 {
		args = append(args, "-bepath="+s.BEPath)
	}

	if len(s.ProfileName) > 0 {
		args = append(args, "-name="+s.ProfileName)
	}

	c := exec.Command(armaExecutable, args...)
	processCh := make(chan *os.Process)
	go func() {
		processCh <- c.Process
	}()
	c.Start()

	return processCh
}
