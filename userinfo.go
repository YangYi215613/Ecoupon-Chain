package main

type UserInfo struct {
	Balance uint64  // 用户余额
	Count uint64  // 交易次数
	CouponList []*Coupon  // 用户卡券列表
}


// 区块要进行hash，所以必须传值
func NewUserInfo() *UserInfo{
	// 初始化用户拥有20个卡券数据
	var couponList []*Coupon

	for i:=0; i<20; i++ {
		couponList = append(couponList, NewCoupon())
	}

	return &UserInfo{
		Balance: 10000,
		Count: 0,
		CouponList: couponList,
	}
}