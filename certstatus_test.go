package main

import (
	"bytes"
	"errors"
	"golang.org/x/crypto/ocsp"
	"io/ioutil"
	"testing"
)

type stubHTTPClient struct{}

func (stubHTTPClient) get(url string) ([]byte, error) {
	if url == "http://cacerts.digicert.com/DigiCertSHA2SecureServerCA.crt" {
		return ioutil.ReadFile("./testdata/issuer.pem")
	}

	return nil, errors.New("Unrecognised URL: " + url)
}

func TestGetIssuerCert(t *testing.T) {
	cert, err := readCertificate("./testdata/certificate.pem")
	if err != nil {
		t.Fatal(err)
	}

	client := &stubHTTPClient{}
	issCert, err := getIssuerCertificate(client, cert)
	if err != nil {
		t.Fatal(err)
	}

	if issCert.Issuer.CommonName != "DigiCert Global Root CA" {
		t.Fatal(issCert.Issuer.CommonName)
	}
}

func TestReadCertificate(t *testing.T) {
	_, err := readCertificate("./testdata/certificate.pem")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetOCSPServer(t *testing.T) {
	cert, _ := readCertificate("./testdata/certificate.pem")
	server, err := getOCSPServer(cert)
	if server != "http://ocsp.digicert.com" {
		t.Fatal(err)
	}
}

func TestCertificateFromBytesNoCertificate(t *testing.T) {
	in, _ := ioutil.ReadFile("./testdata/private_key.pem")
	_, err := certificateFromBytes(in)
	if err == nil {
		t.Fatal("should return error")
	}
}

func TestPrintStatusResponse(t *testing.T) {
	ocsp_der, _ := ioutil.ReadFile("./testdata/twitter_ocsp_response_v1.der")
	resp, err := ocsp.ParseResponse(ocsp_der, nil)

	if err != nil {
		t.Fatal(err)
	}

	out = new(bytes.Buffer) // capture output

	expected := "Serial number: 16190166165489431910151563605275097819\n\n" +
		"Status: Good\n\n" +
		"Produced at: 2017-12-23 06:30:33 +0000 UTC\n" +
		"This update: 2017-12-23 06:30:33 +0000 UTC\n" +
		"Next update: 2017-12-30 05:45:33 +0000 UTC\n"

	printStatusResponse(resp)

	got := out.(*bytes.Buffer).String()
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestStatusMessage(t *testing.T) {
	status := statusMessage(ocsp.Good)
	expected := "Good"

	if status != expected {
		t.Errorf("expected %q, got %q", expected, status)
	}
}

func TestRevocationReason(t *testing.T) {
	reason := revocationReason(ocsp.KeyCompromise)
	expected := "Key compromise"

	if reason != expected {
		t.Errorf("expected %q, got %q", expected, reason)
	}
}
