<script setup>

</script>

<template>
<div class="loader"></div>
</template>

<style scoped>
/* HTML: <div class="loader"></div> */
.loader {
  width: 50px;
  aspect-ratio: 1;
  display: grid;
  --c:#CFF09E 25%,#79BD9A 0 50%,#0B486B 0 75%,#CFF09E 0;
}
.loader:before,
.loader:after {
  content: "";
  background:
    linear-gradient( 90deg,var(--c)),
    linear-gradient(-90deg,var(--c));
  background-size: 400% 52%;
  background-repeat: no-repeat;
  animation: l13 2s infinite; 
}
@keyframes l13 {
  0%       {background-position:calc(3*100%/3) 0,calc(0*100%/3) 100%}
  23%,33%  {background-position:calc(2*100%/3) 0,calc(1*100%/3) 100%}
  56%,66%  {background-position:calc(1*100%/3) 0,calc(2*100%/3) 100%}
  90%,100% {background-position:calc(0*100%/3) 0,calc(3*100%/3) 100%}
}
</style>