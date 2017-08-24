package acceptance_test

import (
	"time"

	acceptance "github.com/cloudfoundry/bosh-bootloader/acceptance-tests"
	"github.com/cloudfoundry/bosh-bootloader/acceptance-tests/actors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("rotate ssh key test", func() {
	var (
		bbl           actors.BBL
		configuration acceptance.Config
	)

	BeforeEach(func() {
		var err error
		configuration, err = acceptance.LoadConfig()
		Expect(err).NotTo(HaveOccurred())

		bbl = actors.NewBBL(configuration.StateFileDir, pathToBBL, configuration, "rotate-env")
	})

	AfterEach(func() {
		session := bbl.Destroy()
		Eventually(session, 10*time.Minute).Should(gexec.Exit())
	})

	It("is able to bbl up without credhub and rotate the director's ssh key", func() {
		session := bbl.Up("--name", bbl.PredefinedEnvID())
		Eventually(session, 40*time.Minute).Should(gexec.Exit(0))
		sshKey := bbl.SSHKey()
		Expect(sshKey).NotTo(BeEmpty())

		By("checking if ssh'ing works", func() {
			err := sshToDirector(bbl, "jumpbox")
			Expect(err).NotTo(HaveOccurred())
		})

		session = bbl.Rotate()
		Eventually(session, 40*time.Minute).Should(gexec.Exit(0))

		rotatedKey := bbl.SSHKey()
		Expect(rotatedKey).NotTo(BeEmpty())
		Expect(rotatedKey).NotTo(Equal(sshKey))

		By("checking that ssh still works", func() {
			err := sshToDirector(bbl, "jumpbox")
			Expect(err).NotTo(HaveOccurred())
		})
	})

	It("is able to bbl up with credhub and rotate the jumpbox's ssh key", func() {
		session := bbl.Up("--credhub", "--name", bbl.PredefinedEnvID())
		Eventually(session, 40*time.Minute).Should(gexec.Exit(0))
		sshKey := bbl.SSHKey()
		Expect(sshKey).NotTo(BeEmpty())

		session = bbl.Rotate()
		Eventually(session, 40*time.Minute).Should(gexec.Exit(0))

		rotatedKey := bbl.SSHKey()
		Expect(rotatedKey).NotTo(BeEmpty())
		Expect(rotatedKey).NotTo(Equal(sshKey))
	})
})
