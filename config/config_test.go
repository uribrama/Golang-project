package config

import (
	"fmt"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	environmentEnv = "GO_ENV"
)

var envVars = map[string]string{
	environmentEnv: "test",
}

var expectedEnvVar = EnvVar{
	Environment: envVars[environmentEnv],
}

var _ = Describe("EnvVar", func() {
	var sut *configProvider

	BeforeEach(func() {
		sut = &configProvider{}
	})

	Describe("parseEnvironmentVariables", func() {
		BeforeEach(func() {
			for envVar, expectedVal := range envVars {
				Ω(os.Setenv(envVar, expectedVal)).Should(Succeed())
			}
		})

		AfterEach(func() {
			for envVar := range envVars {
				Ω(os.Unsetenv(envVar)).Should(Succeed())
			}
		})

		Context("when all expected environment variables are set", func() {
			It("returns EnvVar", func() {
				sut.parseEnvironmentVariables()
				Ω(sut.proyectConfig.EnvVar.Environment).Should(Equal(expectedEnvVar.Environment))
			})
		})

		Context("when expected environment variables are not set", func() {
			UnsetAndAssert := func(env string) {
				BeforeEach(func() {
					Ω(os.Unsetenv(env)).Should(Succeed())
				})
				It("should error appropriately", func() {
					Ω(sut.parseEnvironmentVariables).
						Should(PanicWith(MatchError(ContainSubstring(fmt.Sprintf(`"%s" is not set`, env)))))
				})
			}

			for env := range envVars {
				func(env string) {
					Context("with "+env+" missing", func() {
						UnsetAndAssert(env)
					})
				}(env)
			}
		})
	})
})
