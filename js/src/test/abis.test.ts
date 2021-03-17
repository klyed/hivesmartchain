import * as assert from 'assert';
import {hsc, compile} from "../test";

describe('Abi', function () {

  const source = `
pragma solidity >=0.0.0;

contract random {
	function getRandomNumber() public pure returns (uint) {
		return 55;
	}
}
  `
  // TODO: understand why abi storage isn't working
  it('Call contract via hsc side Abi', async () => {
    const {abi, code} = compile(source, 'random')
    const contractIn: any = await hsc.contracts.deploy(abi, code)
    await hsc.namereg.set('random', contractIn.address)
    const entry = await hsc.namereg.get('random')
    const address = entry.getData();
    console.log(address)
    const contractOut: any = await hsc.contracts.fromAddress(address)
    const number = await contractOut.getRandomNumber()
    assert.strictEqual(number[0], 55)
  })
})
