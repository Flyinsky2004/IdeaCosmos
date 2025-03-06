Component({
  data: {
    selected: 0,
    list: [
      {
        pagePath: "/pages/index/index",
        text: "浏览",
        icon: "success"
      },
      {
        pagePath: "/pages/notification/index",
        text: "通知",
        icon: "info"
      },
      {
        pagePath: "/pages/message/index",
        text: "私信",
        icon: "chat"
      },
      {
        pagePath: "/pages/profile/index",
        text: "我的",
        icon: "personal"
      }
    ]
  },
  methods: {
    switchTab(e) {
      const data = e.currentTarget.dataset;
      const url = data.path;
      wx.switchTab({
        url
      });
      this.setData({
        selected: data.index
      });
    }
  }
}); 