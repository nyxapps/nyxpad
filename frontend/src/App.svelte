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
  let backgroundColor: string = $state("#ffffff");

  let showFontSize: boolean = $state(false);
  let timeoutId: number | null = null;

  ConfigGet("font-family").then((value) => {
    fontFamily = value;
  });

  ConfigGet("font-size").then((value) => {
    fontSize = value;
  });

  ConfigGet("background-color").then((value) => {
    backgroundColor = value;
  });

  //? FILE FUNCTIONS ?//
  function newFile() {
    file = "";
    content = "";
  }

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
    ReadFile()
      .then((value) => {
        file = value.file;
        content = value.content;
      })
      .catch((error) => console.error(error));
  }

  function handleKeyDown(e: KeyboardEvent) {
    // handle tab
    if (e.key === "Tab" && !e.ctrlKey && !e.metaKey && !e.altKey) {
      e.preventDefault();
      const start = textAreaElement.selectionStart;
      const end = textAreaElement.selectionEnd;
      content = content.substring(0, start) + "\t" + content.substring(end);
      textAreaElement.selectionStart = textAreaElement.selectionEnd = start + 1;
    }
  }

  function showFontSizeIndicator() {
    showFontSize = true;

    // Clear any existing timeout
    if (timeoutId) {
      clearTimeout(timeoutId);
    }

    // Set new timeout
    timeoutId = setTimeout(() => {
      showFontSize = false;
      timeoutId = null;
    }, 2000);
  }

  onMount(() => {
    EventsOn("newRequested", newFile);
    EventsOn("saveRequested", saveFile);
    EventsOn("openRequested", openFile);

    const keyDownEvent = (e: KeyboardEvent) => {
      if (e.ctrlKey || e.metaKey) {
        const currentSize = parseInt(fontSize.replace("px", ""));

        // Handle zoom in (Ctrl/Cmd + Plus)
        if (e.key === "+" || e.key === "=") {
          // "=" is the unshifted key for "+"
          e.preventDefault();
          console.log(currentSize);
          fontSize = `${currentSize + 1}px`;
          showFontSizeIndicator();
        }

        // Handle zoom out (Ctrl/Cmd + Minus)
        if (e.key === "-" || e.key === "_") {
          // "_" is shifted "-"
          e.preventDefault();
          console.log(currentSize);
          fontSize = `${Math.max(8, currentSize - 1)}px`; // Minimum 8px
          showFontSizeIndicator();
        }
      }
    };

    window.addEventListener("keydown", keyDownEvent);
    return () => {
      window.removeEventListener("keydown", keyDownEvent);
    };
  });
</script>

<main class="h-full" style={`background-color:${backgroundColor};`}>
  <textarea
    id="content"
    bind:this={textAreaElement}
    bind:value={content}
    onkeydown={handleKeyDown}
    class={`w-full h-[99%] focus:outline-0 whitespace-pre-wrap`}
    style={`font-family:${fontFamily}; tab-size: 4;font-size:${fontSize};`}
  ></textarea>

  <p
    class="fixed bottom-0 left-0 bg-black text-white duration-300"
    style={`opacity:${showFontSize ? 1 : 0};`}
  >
    {fontSize}
  </p>
</main>
