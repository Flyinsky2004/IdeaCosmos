<template>
  <div class="flow-processor">
    <div class="inputs">
      <h2>Inputs</h2>
      <ul>
        <li v-for="framework in frameworks" :key="framework">
          <div class="framework">{{ framework }}</div>
        </li>
      </ul>
    </div>
    <div class="processor">
      <h2>Processor</h2>
      <div class="process-box">
        <p>Processing Data...</p>
      </div>
    </div>
    <div class="outputs">
      <h2>Outputs</h2>
      <div class="output">
        <div class="code-block">
          <h3>HTML</h3>
          <pre>{{ htmlOutput }}</pre>
        </div>
        <div class="code-block">
          <h3>CSS</h3>
          <pre>{{ cssOutput }}</pre>
        </div>
        <div class="code-block">
          <h3>JavaScript</h3>
          <pre>{{ jsOutput }}</pre>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'FlowProcessor',
  data() {
    return {
      frameworks: ['Vue', 'React', 'Svelte', 'Angular'],
      htmlOutput: '<div>Hello World</div>',
      cssOutput: 'div { color: red; }',
      jsOutput: 'console.log("Hello World");',
    };
  },
};
</script>

<style scoped>
.flow-processor {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 20px;
}

.inputs, .outputs, .processor {
  flex: 1;
}

.inputs ul {
  list-style-type: none;
  padding: 0;
}

.inputs .framework {
  padding: 10px;
  margin: 5px 0;
  background-color: #f5f5f5;
  border: 1px solid #ddd;
}

.processor .process-box {
  background-color: #eef;
  border: 1px solid #ccd;
  padding: 20px;
  text-align: center;
}

.outputs .code-block {
  margin: 10px 0;
  padding: 10px;
  border: 1px solid #ccc;
  background-color: #f9f9f9;
}

.outputs .code-block h3 {
  margin: 0 0 10px;
}
</style>
