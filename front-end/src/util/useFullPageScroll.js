import { ref, onMounted, onUnmounted } from 'vue';

export const useFullPageScroll = () => {
    const currentSection = ref(0);
    const isScrolling = ref(false);
    const sections = ref([]);

    const scrollToSection = (index) => {
        if (index >= 0 && index < sections.value.length) {
            isScrolling.value = true;
            currentSection.value = index;
            
            sections.value[index].scrollIntoView({
                behavior: 'smooth'
            });

            // 等待滚动动画完成
            setTimeout(() => {
                isScrolling.value = false;
            }, 1000); // 滚动动画大约需要1秒
        }
    };

    const handleWheel = (event) => {
        if (isScrolling.value) {
            event.preventDefault();
            return;
        }

        const direction = event.deltaY > 0 ? 1 : -1;
        const nextSection = currentSection.value + direction;

        // 确保在有效范围内
        if (nextSection >= 0 && nextSection < sections.value.length) {
            scrollToSection(nextSection);
        }

        event.preventDefault();
    };

    // 监听滚动结束
    const handleScrollEnd = () => {
        if (!isScrolling.value) return;

        const currentScroll = window.scrollY;
        const windowHeight = window.innerHeight;
        const sectionIndex = Math.round(currentScroll / windowHeight);

        if (sectionIndex !== currentSection.value) {
            currentSection.value = sectionIndex;
        }
    };

    onMounted(() => {
        // 初始化sections
        sections.value = Array.from(document.querySelectorAll('.h-screen'));
        
        // 确保从顶部开始
        window.scrollTo(0, 0);
        currentSection.value = 0;

        // 添加事件监听
        window.addEventListener('wheel', handleWheel, { passive: false });
        window.addEventListener('scrollend', handleScrollEnd);

        // 添加触摸事件支持
        let touchStartY = 0;
        window.addEventListener('touchstart', (e) => {
            touchStartY = e.touches[0].clientY;
        });

        window.addEventListener('touchend', (e) => {
            if (isScrolling.value) return;
            
            const touchEndY = e.changedTouches[0].clientY;
            const direction = touchStartY > touchEndY ? 1 : -1;
            const nextSection = currentSection.value + direction;

            if (Math.abs(touchStartY - touchEndY) > 50 && // 确保是有意义的滑动
                nextSection >= 0 && 
                nextSection < sections.value.length) {
                scrollToSection(nextSection);
            }
        });
    });

    onUnmounted(() => {
        window.removeEventListener('wheel', handleWheel);
        window.removeEventListener('scrollend', handleScrollEnd);
    });

    return {
        currentSection,
        scrollToSection
    };
};