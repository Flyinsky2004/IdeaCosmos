import { ref, onMounted, onUnmounted } from 'vue';

export const useFullPageScroll = () => {
    const currentSection = ref(0);
    const isScrolling = ref(false);
    const sections = ref([]);

    const initializeSections = () => {
        sections.value = Array.from(document.querySelectorAll('.h-screen'));
        window.scrollTo(0, 0);
    };

    const handleWheel = (event) => {
        if (isScrolling.value) {
            event.preventDefault();
            return;
        }

        const direction = event.deltaY > 0 ? 1 : -1;
        const nextSection = Math.min(
            Math.max(currentSection.value + direction, 0),
            sections.value.length - 1
        );

        if (nextSection !== currentSection.value) {
            isScrolling.value = true;
            currentSection.value = nextSection;

            sections.value[nextSection].scrollIntoView({
                behavior: 'smooth'
            });
            isScrolling.value = false;
        }

        event.preventDefault();
    };

    onMounted(() => {
        initializeSections();
        window.addEventListener('wheel', handleWheel, { passive: false });
    });

    onUnmounted(() => {
        window.removeEventListener('wheel', handleWheel);
    });

    return {
        currentSection
    };
};