<script lang="ts">
  import { onMount } from "svelte";
  import { Flashlight, FlashlightOff, Info, Loader2 } from "lucide-svelte";
  import QrScanner from "qr-scanner";
  import { goto } from "$app/navigation";

  export let loading: boolean = false;
  export let error: string | undefined;

  let videoRef: HTMLVideoElement | null = null;
  let qrScanner: QrScanner | null = null;
  let torchTrack: MediaStreamTrack | null = null;

  let hasCamera = true;
  let cameraError: string = "";
  let torchSupported = false;
  let torchEnabled = false;
  let torchError: string | null = null;

  import { scannedEquipmentId } from "$lib/stores/equipmentContext";

  const setupScanner = async () => {
    try {
      const devices = await navigator.mediaDevices.enumerateDevices();
      const videoDevices = devices.filter((d) => d.kind === "videoinput");
      if (videoDevices.length === 0)
        throw new Error("No video input devices found.");

      const selectedDevice = videoDevices[videoDevices.length - 1];
      const stream = await navigator.mediaDevices.getUserMedia({
        video: { deviceId: { exact: selectedDevice.deviceId } },
      });

      const track = stream.getVideoTracks()[0];
      const capabilities = track.getCapabilities?.() ?? {};
      torchSupported = "torch" in capabilities;
      torchTrack = track;

      if (videoRef) {
        videoRef.srcObject = stream;
        videoRef.onloadedmetadata = () => videoRef?.play().catch(console.error);
      }

      qrScanner = new QrScanner(
        videoRef!,
        // @ts-expect-error: result.data is valid
        (result) => {
          const id = result.data;
          if (id) {
            goto("/equipment?scanned=" + id);
          }
        },
        {
          highlightScanRegion: true,
          highlightCodeOutline: true,
          deviceId: selectedDevice.deviceId,
        }
      );

      await qrScanner.start();

      hasCamera = true;
      cameraError = "";
    } catch (err) {
      console.error(err);
      hasCamera = false;
      cameraError =
        "Unable to access camera. Please ensure permissions are granted.";
    }
  };

  const toggleTorch = async () => {
    if (!torchTrack) {
      torchError = "No torch-capable track available.";
      return;
    }

    try {
      const enableTorch = !torchEnabled;
      await torchTrack.applyConstraints({
        advanced: [{ torch: enableTorch } as MediaTrackConstraintSet],
      });
      torchEnabled = enableTorch;
      torchError = null;
    } catch (err) {
      console.error(err);
      torchError = "Failed to toggle torch.";
    }
  };

  onMount(() => {
    setupScanner();

    return () => {
      qrScanner?.destroy();
      torchTrack?.stop();

      const stream = videoRef?.srcObject;
      if (stream instanceof MediaStream) {
        stream.getTracks().forEach((t: MediaStreamTrack) => t.stop());
      }
    };
  });
</script>

<div class="fixed inset-0 bg-black">
  <div class="relative w-full h-full">
    <video
      bind:this={videoRef}
      class="w-full h-full object-cover"
      autoplay
      muted
      playsinline
    ></video>

    <div
      class="absolute top-0 left-0 right-0 p-6 z-10 flex items-center justify-between"
    >
      <div class="bg-black/60 backdrop-blur-sm rounded-full px-4 py-2">
        <div class="flex items-center space-x-2">
          <div class="w-2 h-2 bg-green-400 rounded-full animate-pulse"></div>
          <span class="text-white text-sm font-medium"
            >Scanning for QR codes...</span
          >
        </div>
      </div>

      {#if torchSupported}
        <button
          on:click={toggleTorch}
          class={`p-3 rounded-full backdrop-blur-sm transition-all ${
            torchEnabled
              ? "bg-yellow-500/80 text-white"
              : "bg-black/60 text-white hover:bg-black/80"
          }`}
        >
          {#if torchEnabled}
            <Flashlight class="w-6 h-6" />
          {:else}
            <FlashlightOff class="w-6 h-6" />
          {/if}
        </button>
      {/if}
    </div>

    <div class="absolute bottom-0 left-0 right-0 p-6 z-10">
      <div class="text-center">
        <div
          class="bg-black/60 backdrop-blur-sm rounded-2xl px-6 py-4 inline-block"
        >
          <h2 class="text-white text-xl font-semibold mb-2">
            Scan Equipment QR Code
          </h2>
          <p class="text-white/80 text-sm">
            Position the QR code within the frame to scan equipment information
          </p>
        </div>
      </div>
    </div>

    {#if loading}
      <div
        class="absolute inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center z-20"
      >
        <div class="bg-white rounded-2xl p-8 text-center max-w-sm mx-4">
          <Loader2 class="w-12 h-12 text-blue-600 animate-spin mx-auto mb-4" />
          <h3 class="text-lg font-semibold text-gray-900 mb-2">
            Loading Equipment
          </h3>
          <p class="text-gray-600">Fetching equipment information...</p>
        </div>
      </div>
    {/if}

    {#if error || cameraError}
      <div
        class="absolute inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center z-20"
      >
        <div class="bg-white rounded-2xl p-8 text-center max-w-sm mx-4">
          <div
            class="w-16 h-16 bg-red-100 rounded-full flex items-center justify-center mx-auto mb-4"
          >
            <Info class="w-8 h-8 text-red-600" />
          </div>
          <h3 class="text-lg font-semibold text-gray-900 mb-2">
            {cameraError ? "Camera Error" : "Equipment Not Found"}
          </h3>
          <p class="text-gray-600 mb-6">{cameraError || error}</p>
          <button
            on:click={() => window.location.reload()}
            class="w-full bg-blue-600 text-white py-3 px-4 rounded-lg hover:bg-blue-700 transition-colors font-medium"
          >
            Try Again
          </button>
        </div>
      </div>
    {/if}

    {#if torchError}
      <div class="absolute top-20 left-0 right-0 px-6 z-30 flex justify-center">
        <div
          class="bg-red-600 text-white text-sm rounded-lg px-4 py-2 shadow-lg"
        >
          {torchError}
        </div>
      </div>
    {/if}
  </div>
</div>
