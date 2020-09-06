package iotmakerDocker

import (
	"bytes"
	"errors"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"github.com/helmutkemper/iotmaker.docker/util"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

func ExampleDockerSystem_ContainerCreateChangeExposedPortAndStart() {
	var err error
	var containerId string
	var networkId string
	var dockerSys *DockerSystem

	var networkAutoConfiguration *NextNetworkAutoConfiguration
	var networkNextAddress *network.NetworkingConfig

	// English: make a channel to end goroutine
	// Português: monta um canal para terminar a goroutine
	var chProcessEnd = make(chan bool, 1)

	// English: make a channel [optional] to print build output
	// Português: monta o canal [opcional] para imprimir a saída do build
	var chStatus = make(chan ContainerPullStatusSendToChannel, 1)

	// English: make a thread to monitoring and print channel data
	// Português: monta uma thread para imprimir os dados do canal
	go func(chStatus chan ContainerPullStatusSendToChannel, chProcessEnd chan bool) {

		for {
			select {
			case <-chProcessEnd:
				return

			case status := <-chStatus:
				// English: remove this comment to see all build status
				// Português: remova este comentário para vê todo o status da criação da imagem
				//fmt.Printf("image pull status: %+v\n", status)

				if status.Closed == true {
					// fmt.Println("image pull complete!")

					// English: Eliminate this goroutine after process end
					// Português: Elimina a goroutine após o fim do processo
					// return
				}
			}
		}

	}(chStatus, chProcessEnd)

	// English: searches for the folder containing the test server
	// Português: procura pela pasta contendo o servidor de teste
	var smallServerPath string
	smallServerPath, err = util.FileFindRecursively("small_test_server_port_3000")
	if err != nil {
		panic(err)
	}

	// English: turns the path into an absolute path
	// Português: transforma o caminho em caminho absoluto
	smallServerPath, err = filepath.Abs(smallServerPath)
	if err != nil {
		panic(err)
	}

	// English: 'static' folder path
	// Português: caminho da pasta 'static'
	var smallServerPathStatic string
	smallServerPathStatic = smallServerPath + string(os.PathSeparator) + "static"

	// English: 'index.html' file path
	// Português: caminho do arquivo 'index.html'
	var smallServerPathStaticContent string
	smallServerPathStaticContent = smallServerPathStatic + string(os.PathSeparator) + "index.html"

	// English: create a new default client. Please, use: err, dockerSys = factoryDocker.NewClient()
	// Português: cria um novo cliente com configurações padrão. Por favor, usr: err, dockerSys = factoryDocker.NewClient()
	dockerSys = &DockerSystem{}
	dockerSys.ContextCreate()
	err = dockerSys.ClientCreate()
	if err != nil {
		panic(err)
	}

	// English: garbage collector and deletes networks and images whose name contains the term 'delete'
	// Português: coletor de lixo e apaga redes e imagens cujo o nome contém o temo 'delete'
	err = dockerSys.RemoveAllByNameContains("delete")
	if err != nil {
		panic(err)
	}

	// English: create a network named 'network_delete_before_test'
	// Português: cria uma nova rede de nome 'network_delete_before_test'
	networkId, networkAutoConfiguration, err = dockerSys.NetworkCreate(
		"network_delete_before_test",
		KNetworkDriveBridge,
		"local",
		"10.0.0.0/16",
		"10.0.0.1",
	)
	if err != nil {
		panic(err)
	}

	if networkId == "" {
		err = errors.New("network id was not generated")
		panic(err)
	}

	// English: get next ip address from network, '10.0.0.2'
	// Português: pega o próxima endereço da rede, '10.0.0.2'
	err, networkNextAddress = networkAutoConfiguration.GetNext()
	if err != nil {
		panic(err)
	}

	// English: build a new image from folder 'small_test_server_port_3000'
	// Português: monta uma imagem a partir da pasta 'small_test_server_port_3000'
	err = dockerSys.ImageBuildFromFolder(
		smallServerPath,
		[]string{
			"image_server_delete_before_test:latest", // image name
		},
		&chStatus, // [channel|nil]
	)
	if err != nil {
		panic(err)
	}

	// English: mount and start a container
	// Português: monta i inicializa o container
	containerId, err = dockerSys.ContainerCreateChangeExposedPortAndStart(
		"image_server_delete_before_test:latest", // image name
		"container_delete_before_test",           // container name
		KRestartPolicyUnlessStopped,              // restart policy
		[]mount.Mount{ // mount volumes
			{
				Type: KVolumeMountTypeBindString, // bind - is the type for mounting host dir
				// (real folder inside computer where this
				// code work)

				Source: smallServerPathStatic, // path inside host machine
				Target: "/static",             // path inside image
			},
		},
		networkNextAddress, // [optional] container network
		[]nat.Port{ // image port list
			"3000/tcp", // small server port :3000
		},
		[]nat.Port{ // host computer new port list
			"9002/tcp", // http://localhost:9002
		},
	)
	if err != nil {
		panic(err)
	}

	if containerId == "" {
		err = errors.New("container id was not generated")
		panic(err)
	}

	// English: ends a goroutine
	// Português: termina a goroutine
	chProcessEnd <- true

	var page []byte
	var found int
	var response *http.Response

	// English: contents of the static server index.html page
	// Português: conteúdo da página index.html do servidor estático
	var htmlIndexContent1 = []byte("<!DOCTYPE html><html><body><p>I'm alive!</p></body></html>")

	err = ioutil.WriteFile(smallServerPathStaticContent, htmlIndexContent1, os.ModePerm)
	if err != nil {
		panic(err)
	}

	// English: takes the content from the server for testing
	// Português: pega o conteúdo do servidor para teste
	response, err = http.Get("http://localhost:9002")
	if err != nil {
		panic(err)
	}

	page, err = ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	// English: compare content
	// Português: compara o conteúdo
	found = bytes.Compare(page, htmlIndexContent1)
	if found == 1 {
		err = errors.New("small server content error")
		panic(err)
	}

	// English: contents of the static server index.html page
	// Português: conteúdo da página index.html do servidor estático
	htmlIndexContent1 = []byte("<!DOCTYPE html><html><body><p>I'm alive Dr. Frankenstein!</p></body></html>")

	err = ioutil.WriteFile(smallServerPathStaticContent, htmlIndexContent1, os.ModePerm)
	if err != nil {
		panic(err)
	}

	// English: takes the content from the server for testing
	// Português: pega o conteúdo do servidor para teste
	response, err = http.Get("http://localhost:9002")
	if err != nil {
		panic(err)
	}

	page, err = ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	// English: compare content
	// Português: compara o conteúdo
	found = bytes.Compare(page, htmlIndexContent1)
	if found == 1 {
		err = errors.New("small server content error")
		panic(err)
	}

	// English: garbage collector and deletes networks and images whose name contains the term 'delete'
	// Português: coletor de lixo e apaga redes e imagens cujo o nome contém o temo 'delete'
	err = dockerSys.RemoveAllByNameContains("delete")
	if err != nil {
		panic(err)
	}

	// Output:
	//
}
