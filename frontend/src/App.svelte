<script lang="ts">
  import { onMount } from "svelte";
  import {
    ConfigGet,
    ReadFile,
    SaveFile,
    WriteFile,
  } from "../wailsjs/go/main/App.js";
  import { EventsOn } from "../wailsjs/runtime/runtime.js";

  let textAreaElement: HTMLTextAreaElement;
  let content: string = $state("");
  let file: string | undefined = undefined;

  let fontFamily: string = $state("sans-serif");
  let fontSize: string = $state("16px");

  ConfigGet("font-family").then((value) => {
    fontFamily = value;
  });

  ConfigGet("font-size").then((value) => {
    fontSize = value;
  });

  // ? FILE FUNCTIONS ?//
  function saveFile() {
    if (!file) {
      SaveFile(content)
        .then((value) => {
          file = value;
        })
        .catch((error) => console.error(error));
    } else {
      WriteFile(file, content).catch((error) => console.error(error));
    }
  }

  function openFile() {
    ReadFile().then((value) => {
      file = value.file;
      content = value.content;
    });
  }

  function handleSelectApp(e: KeyboardEvent) {
    if ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === "a") {
      e.preventDefault();
      textAreaElement.select();
    }
  }

  onMount(() => {
    EventsOn("saveRequested", saveFile);
    EventsOn("openRequested", openFile);
  });
</script>

<main class="h-full">
  <textarea
    name="content"
    id="content"
    bind:this={textAreaElement}
    bind:value={content}
    onkeydown={handleSelectApp}
    class={`w-full h-[99%] resize-none focus:outline-0 text-[${fontSize}]`}
    style={`font-family:${fontFamily};`}
  ></textarea>
</main>
