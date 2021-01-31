package persistance

// PacketExist return if packet exist
//
// If more than one parameter is not NIL all must match
func PacketExist(bucket string, opts PacketFindOpts) (found bool, packetInfo PacketInfo, err error) {
	found, _, _, err = packetGet(bucket, opts)
	return
}
