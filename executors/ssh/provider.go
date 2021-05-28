package ssh

import (
	"io"
	"io/ioutil"

	"github.com/projecteru2/phistage/common"
	"github.com/projecteru2/phistage/executors"
	"github.com/projecteru2/phistage/store"

	"golang.org/x/crypto/ssh"
)

type SSHJobExecutorProvider struct {
	config *common.Config
	store  store.Store
}

func NewSSHJobExecutorProvider(config *common.Config, store store.Store) (*SSHJobExecutorProvider, error) {
	return &SSHJobExecutorProvider{
		config: config,
		store:  store,
	}, nil
}

func (s *SSHJobExecutorProvider) GetName() string {
	return "ssh"
}

func (s *SSHJobExecutorProvider) GetJobExecutor(job *common.Job, phistage *common.Phistage, output io.Writer) (executors.JobExecutor, error) {
	key, err := ioutil.ReadFile(s.config.SSH.PrivateKey)
	if err != nil {
		return nil, err
	}

	singer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}

	client, err := ssh.Dial("tcp", s.config.SSH.Address, &ssh.ClientConfig{
		User: s.config.SSH.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(singer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		return nil, err
	}
	return NewSSHJobExecutor(job, phistage, output, client, s.store, s.config)
}
