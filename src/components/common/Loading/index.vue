<template>
  <section class="loaidng-dots-container">
    <div class="loaidng-dot" />
    <div class="loaidng-dot" />
    <div class="loaidng-dot" />
    <div class="loaidng-dot" />
    <div class="loaidng-dot" />
  </section>
</template>

<style>
/* From Uiverse.io by adamgiebl */
.loaidng-dots-container {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  width: 100%;
}

.loaidng-dot {
  height: 10px;
  width: 10px;
  margin-right: 10px;
  border-radius: 10px;
  background-color: #b3d4fc;
  animation: pulse 1.5s infinite ease-in-out;
}

.loaidng-dot:last-child {
  margin-right: 0;
}

.loaidng-dot:nth-child(1) {
  animation-delay: -0.3s;
}

.loaidng-dot:nth-child(2) {
  animation-delay: -0.1s;
}

.loaidng-dot:nth-child(3) {
  animation-delay: 0.1s;
}

@keyframes pulse {
  0% {
    transform: scale(0.8);
    background-color: #b3d4fc;
    box-shadow: 0 0 0 0 rgba(178, 212, 252, 0.7);
  }

  50% {
    transform: scale(1.2);
    background-color: #6793fb;
    box-shadow: 0 0 0 10px rgba(178, 212, 252, 0);
  }

  100% {
    transform: scale(0.8);
    background-color: #b3d4fc;
    box-shadow: 0 0 0 0 rgba(178, 212, 252, 0.7);
  }
}
</style>
