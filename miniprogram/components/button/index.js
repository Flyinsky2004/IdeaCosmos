Component({
  /**
   * 组件的属性列表
   */
  properties: {
    type: {
      type: String,
      value: 'primary' // primary, success, danger
    },
    disabled: {
      type: Boolean,
      value: false
    },
    block: {
      type: Boolean,
      value: true
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
      if (this.data.disabled) return;
      this.triggerEvent('tap');
    }
  }
}) 