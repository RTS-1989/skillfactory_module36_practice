<template>
  <div>
    <h2 class="mx-5 my-5">GoNews - агрегатор новостей.</h2>
    <div v-for="post in news" :key="post.id">
      <v-card elevation="10" outlined class="mx-5 my-5">
        <v-card-title>
          <a :href="post.link" target="_blank"> {{ post.title }} </a>
        </v-card-title>
        <v-card-text>
          {{ post.content }}
          <v-card-subtitle>
            {{ new Date(post.pub_time * 1000) }}
          </v-card-subtitle>
        </v-card-text>
      </v-card>
    </div>
  </div>
</template>

<script>
export default {
  name: "News",
  data() {
    return {
      news: [],
    };
  },
  mounted() {
    let url = "http://" + window.location.hostname + ":" + window.location.port + "/news/40";
    fetch(url)
      .then((response) => response.json())
      .then((data) => (this.news = data));
  },
};
</script>

<style scoped>
</style>
