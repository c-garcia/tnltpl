package main

import (
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"
)

const (
	NGROK_TUNNELS = "http://localhost:4040/api/tunnels"
)

type cliArgs struct {
	InFileName  string
	OutFileName string
	TunnelName  string
}

func parseArgs() cliArgs {
	var args cliArgs
	flag.StringVar(&args.InFileName, "in", "", "template to process")
	flag.StringVar(&args.OutFileName, "out", "", "output file")
	flag.StringVar(&args.TunnelName, "name", "command_line", "tunnel name")
	flag.Parse()
	if args.InFileName == "" {
		log.Fatal("You must specify in input file")
	}
	if args.OutFileName == "" {
		log.Fatal("You must specify an output file")
	}
	return args
}

func getTemplate(fileName string) *template.Template {
	tpl, err := template.ParseFiles(fileName)

	if err != nil {
		log.Fatalf("Problem processing input file: %v", err)
	}
	return tpl
}

type tunnel struct {
	Name      string `json:"name"`
	PublicUrl string `json:"public_url"`
}

type tunnelsResp struct {
	Tunnels []tunnel `json:"tunnels"`
}

func findAllTunnels() []tunnel {
	client := &http.Client{}
	req, err := http.NewRequest("GET", NGROK_TUNNELS, nil)
	if err != nil {
		log.Fatalf("Error when creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error when requesting tunnel list: %v", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error when processing response body %v", err)
	}

	var jsonResult tunnelsResp
	if err = json.Unmarshal(body, &jsonResult); err != nil {
		log.Fatalf("Error parsing JSON response %v", err)
	}

	return jsonResult.Tunnels
}

func findTunnelWithName(name string) (tunnel, bool) {
	tunnels := findAllTunnels()
	var res tunnel
	found := false
	for _, t := range tunnels {
		if t.Name == name {
			found = true
			res = t
			break
		}
	}
	return res, found
}

func main() {

	cmdArgs := parseArgs()
	tpl := getTemplate(cmdArgs.InFileName)

	var t tunnel
	var pres bool
	if t, pres = findTunnelWithName(cmdArgs.TunnelName); !pres {
		log.Fatalf("Tunnel with name %s not currently present", cmdArgs.TunnelName)
	}

	var templateData struct {
		Url string
	}
	templateData.Url = t.PublicUrl
	var out io.WriteCloser
	if cmdArgs.OutFileName == "-" {
		out = os.Stdout
	} else {
		var err error
		if out, err = os.Create(cmdArgs.OutFileName); err != nil {
			log.Fatalf("Error when creating output file: %v", err)
		}
	}

	tpl.Execute(out, templateData)

}
