package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/gomail.v2"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"context"
)

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(url string) (outFile *os.File, err error) {

	tmpfile, err := ioutil.TempFile("/tmp/", "ios_")
	if err != nil {
		fmt.Printf("%s", err.Error())
		log.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	outFile, err = os.Create(fmt.Sprintf("%s.jpg", tmpfile.Name()))
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		log.Fatal(err)
	}
	defer outFile.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		log.Fatal(err)
	}

	return
}

type MyEvent struct {
        Name string ""
}

func HandleRequest(ctx context.Context, name MyEvent) (string, error) {
        return fmt.Sprintf("Hello %s!", name.Name ), nil
}

func PhotoMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	m := make(chan string)

	defer r.Body.Close()

	body, _ := ioutil.ReadAll(r.Body)

	go DownloadAndSend(fmt.Sprintf("%s", body), m)
	rec := <-m
	fmt.Printf("%s\n", rec)
}

func sendPhoto(imagePath string, emailAddr string) {
	m := gomail.NewMessage()
	m.SetHeader("From", "graham@XXX.com")
	m.SetHeader("To", emailAddr)
	m.SetHeader("Subject", "New Picture")
	m.SetBody("text/html", "New picture from Apple Shared Photo Album")
	m.Attach(imagePath)

	d := gomail.NewDialer("smtp.gmail.com", 587, "graham@gXXXXX.com", "npaXXXX")

	if err := d.DialAndSend(m); err != nil {
		fmt.Printf("%s\n", err.Error())
		log.Fatal(err)
	}
}

func DownloadAndSend(photoURL string, m chan string) {
	outFile, _ := DownloadFile(photoURL)
	defer os.Remove(outFile.Name())

	sendPhoto(outFile.Name(), "graham@XXX.com")

	m <- "Complete Send File"
}

func main() {
	// lambda.Start(HandleRequest)
//	testSend()
	router := httprouter.New()
	router.PUT("/", PhotoMessage)
	log.Fatal(http.ListenAndServe(":8050", router))
}

func testSend() {
	m := make(chan string)
	testPhoto := "http://cdn4.i-scmp.com/sites/default/files/styles/980x551/public/images/methode/2017/05/19/2b2d8790-3c6a-11e7-8ee3-761f02c18070_1280x720_204107.jpg?itok=k3cLlyz-"
	go DownloadAndSend(testPhoto, m)
	rec := <-m
	fmt.Printf("%s\n", rec)
}
