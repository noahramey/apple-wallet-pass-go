package main

import (
	"archive/zip"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
)

func main() {
	ticketID := "1"
	GeneratePassStructure(ticketID)
	dat, err := RetrieveDataForPass("", ticketID)
	check(err)
	CreatePassJSONFromData(dat, ticketID)
	AssembleManifest()
	ZipPass()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func GeneratePassStructure(ticketID string) {
	// Create the pass folder
	// str := fmt.Sprintf("tmp/%v-pass", ticketID)
	// os.Mkdir(str, 0700)
	// os.Chdir(str)
	// Create the manifest.json file so
	os.Create("./tmp/manifest.json")
	os.Create("./tmp/pass.json")
	os.Link("./assets/tickets/logo.png", "./tmp/logo.png")
	os.Link("./assets/tickets/logo@2x.png", "./tmp/logo@2x.png")
	os.Link("./assets/tickets/icon.png", "./tmp/icon.png")
	os.Link("./assets/tickets/icon@2x.png", "./tmp/icon@2x.png")
}

func RetrieveDataForPass(registrantID string, ticketID string) (map[string]interface{}, error) {
	var err error
	// var r *domain.Registrant
	var data = make(map[string]interface{})
	// sss := svc.db.NewSession()
	// r, err = sss.RegistrantMapper().Find(registrantID, true)
	// rg, err := sss.RegistrationMapper().FindByRegistrantID(registrantID)
	// t := rg.Tickets.Find(ticketID)
	// org for testing purposes
	var org = "fake_company"
	var evn = "event_name"
	var evd = "event_date"
	var evl = "location"
	var td = "1:30 PM"
	var tid = "ticket_id"

	data["formatVersion"] = "1"
	data["passTypeIdentifier"] = "pass.com.fakepass.test"
	data["organizationName"] = org
	data["event_name"] = evn
	data["event_date"] = evd
	data["event_location"] = evl
	data["event_time"] = td
	data["ticket_id"] = tid

	return data, err
}

func CreatePassJSONFromData(data map[string]interface{}, ticketID string) {
	// str := fmt.Sprintf("tmp/%v-pass/pass.json", ticketID)
	str := "./tmp/pass.json"
	// f, err := os.Open(str)
	// check(err)

	// defer f.Close()

	j, err := json.Marshal(data)
	check(err)
	ioutil.WriteFile(str, j, 0700)
	// f.Write(j)
}

func GeneratePassShaHash(filePath string) (string, error) {
	var rs string
	f, err := os.Open(filePath)
	if err != nil {
		return rs, err
	}

	defer f.Close()

	hash := sha1.New()
	_, err = io.Copy(hash, f)
	if err != nil {
		return rs, err
	}

	hashInBytes := hash.Sum(nil)[:20]
	rs = hex.EncodeToString(hashInBytes)

	return rs, nil
}

func AssembleManifest() {
	paths := []string{"./tmp/icon.png", "./tmp/icon@2x.png", "./tmp/logo.png", "./tmp/logo@2x.png", "./tmp/pass.json"}
	hashes := []string{"", "", "", "", ""}
	for i := 0; i < len(paths); i++ {
		h, err := GeneratePassShaHash(paths[i])
		check(err)
		hashes[i] = h
	}

	dat := make(map[string]interface{})
	for i := 0; i < len(paths); i++ {
		dat[paths[i]] = hashes[i]
	}

	j, err := json.Marshal(dat)
	check(err)
	ioutil.WriteFile("./tmp/manifest.json", j, 0700)
}

func SignManifest() {
	f, err := os.Open("./tmp/manifest.json")
	check(err)
	defer f.Close()
}

func ZipPass() {
	filePaths := []string{"./tmp/manifest.json", "./tmp/pass.json", "./tmp/icon.png", "./tmp/icon@2x.png", "./tmp/logo.png", "./tmp/logo@2x.png", "./tmp/pass.json"}

	zipFile, err := os.Create("./tmp/pass.pkpass")
	check(err)
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for _, file := range filePaths {
		zipfile, err := os.Open(file)
		check(err)
		defer zipfile.Close()

		info, err := zipfile.Stat()
		check(err)

		header, err := zip.FileInfoHeader(info)
		check(err)

		header.Name = file
		header.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(header)
		check(err)
		io.Copy(writer, zipfile)
	}
}
