// components/avatar/index.js
Component({
  /**
   * 组件的属性列表
   */
  properties: {
    src: {
      type: String,
      value: ''
    },
    editable: {
      type: Boolean,
      value: false
    },
    overlayText: {
      type: String,
      value: '更换头像'
    }
  },

  /**
   * 组件的初始数据
   */
  data: {

  },

  /**
   * 组件的方法列表
   */
  methods: {
    handleTap() {
      if (!this.data.editable) return;
      this.triggerEvent('tap');
    }
  }
}) 