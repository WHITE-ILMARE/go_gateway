package reverse_proxy

//func NewGrpcLoadBalanceHandler(lb load_balance.LoadBalance) grpc.StreamHandler {
//	return func() grpc.StreamHandler {
//		nextAddr, err := lb.Get("")
//		if err != nil {
//			log.Fatal("get next addr fail")
//		}
//		director := func(ctx context.Context, fullMethodName string) (context.Context, *grpc.ClientConn, error) {
//			c, err := grpc.DialContext(ctx, nextAddr, grpc.WithCodec(proxy.Codec()), grpc.WithInsecure())
//			md, _ := metadata.FromIncomingContext(ctx)
//			outCtx, _ := context.WithCancel(ctx)
//			outCtx = metadata.NewOutgoingContext(outCtx, md.Copy())
//			return outCtx, c, err
//		}
//		return proxy.TransparentHandler(director)
//	}()
//}
