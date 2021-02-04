package bls

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const checkOK = "3e886a4c6e109a151f4105aee65a5192d150ef1fa68d3cd76964a0b086006dbe4324c989deb0e4416c6d6706db1b1910eb2732f08842fb4886067b9ed191109ac2188d76002d2e11da80a3f0ea89fee6b59c834cc478a6bd49cb8a193b1abb16@e96bd0f36b70c5ccc0c4396343bd7d8255b8a526c55fa1e218511fafe6539b8e@04725db195e37aa237cdbbda76270d4a229b6e7a3651104dc58c4349c0388e8546976fe54a04240530b99064e434c90f"
const checkNOK = "2c9a358953f61d34401d7ee4175eec105c476b18baacab371e2f47270035b539d84ad79ba587552b7e38802be00ff7148fc2a9c7a7034ff1e63ee24602ee952235ad14ca7d36e2be617fb2c99ed22a7a2729d86ae9fbb4df06f957ba07fec50e@1e46d9cbb995e30b82485525c29f80ac78aca295a6e88a11c3df8f9a445494bb@be8c460db180d6254c712ead3aa81935bc9be15b919dd45cb152b3dece04762569778c5e70e7af03fa1c66409d4f4711"

func TestBls_VerifyBLS(t *testing.T) {
	t.Parallel()

	b := NewBLS()
	assert.Nil(t, b.VerifyBLS(splitString(t, checkOK)))
	assert.NotNil(t, b.VerifyBLS(splitString(t, checkNOK)))
}

func splitString(t testing.TB, str string) ([]byte, []byte, []byte) {
	split := strings.Split(str, "@")
	pkBuff, err := hex.DecodeString(split[0])
	require.Nil(t, err)

	msgBuff, err := hex.DecodeString(split[1])
	require.Nil(t, err)

	sigBuff, err := hex.DecodeString(split[2])
	require.Nil(t, err)

	return pkBuff, msgBuff, sigBuff
}
