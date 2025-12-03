// snowfall.js - Interactive Snowfall Effect
class InteractiveSnowfall {
    constructor() {
        this.canvas = null;
        this.ctx = null;
        this.snowflakes = [];
        this.mouseX = 0;
        this.mouseY = 0;
        this.tiltX = 0;
        this.tiltY = 0;
        this.maxSnowflakes = 150;
        this.init();
    }

    init() {
        // Tạo canvas
        this.canvas = document.createElement('canvas');
        this.canvas.id = 'snowfall-canvas';
        this.canvas.style.position = 'fixed';
        this.canvas.style.top = '0';
        this.canvas.style.left = '0';
        this.canvas.style.pointerEvents = 'none';
        this.canvas.style.zIndex = '9999';
        document.body.appendChild(this.canvas);

        this.ctx = this.canvas.getContext('2d');
        this.resizeCanvas();

        // Tạo bông tuyết
        for (let i = 0; i < this.maxSnowflakes; i++) {
            this.snowflakes.push(this.createSnowflake());
        }

        // Event listeners
        window.addEventListener('resize', () => this.resizeCanvas());
        window.addEventListener('mousemove', (e) => this.handleMouseMove(e));

        // Gyroscope cho điện thoại
        if (window.DeviceOrientationEvent) {
            window.addEventListener('deviceorientation', (e) => this.handleDeviceOrientation(e));
        }

        // Bắt đầu animation
        this.animate();
    }

    resizeCanvas() {
        this.canvas.width = window.innerWidth;
        this.canvas.height = window.innerHeight;
    }

    createSnowflake() {
        return {
            x: Math.random() * this.canvas.width,
            y: Math.random() * this.canvas.height - this.canvas.height,
            radius: Math.random() * 3 + 1,
            speed: Math.random() * 1 + 0.5,
            wind: Math.random() * 0.5 - 0.25,
            opacity: Math.random() * 0.5 + 0.5
        };
    }

    handleMouseMove(e) {
        // Tính toán hướng gió dựa trên vị trí chuột
        const centerX = this.canvas.width / 2;
        const centerY = this.canvas.height / 2;

        this.mouseX = (e.clientX - centerX) / centerX;
        this.mouseY = (e.clientY - centerY) / centerY;
    }

    handleDeviceOrientation(e) {
        // Xử lý nghiêng điện thoại
        // beta: nghiêng trước sau (-180 đến 180)
        // gamma: nghiêng trái phải (-90 đến 90)
        if (e.gamma !== null && e.beta !== null) {
            this.tiltX = e.gamma / 90; // Chuẩn hóa -1 đến 1
            this.tiltY = (e.beta - 45) / 90; // Chuẩn hóa
        }
    }

    updateSnowflake(flake) {
        // Tính toán wind effect từ chuột hoặc gyroscope
        const windEffect = this.tiltX !== 0 ? this.tiltX * 2 : this.mouseX * 1.5;
        const gravityEffect = this.tiltY !== 0 ? this.tiltY * 0.5 : this.mouseY * 0.3;

        // Cập nhật vị trí
        flake.x += flake.wind + windEffect;
        flake.y += flake.speed + gravityEffect;

        // Hiệu ứng lắc lư tự nhiên
        flake.wind += (Math.random() - 0.5) * 0.1;
        flake.wind = Math.max(-2, Math.min(2, flake.wind));

        // Reset nếu ra khỏi màn hình
        if (flake.y > this.canvas.height) {
            flake.y = -10;
            flake.x = Math.random() * this.canvas.width;
        }

        if (flake.x > this.canvas.width + 10) {
            flake.x = -10;
        } else if (flake.x < -10) {
            flake.x = this.canvas.width + 10;
        }
    }

    drawSnowflake(flake) {
        this.ctx.beginPath();
        this.ctx.arc(flake.x, flake.y, flake.radius, 0, Math.PI * 2);
        this.ctx.fillStyle = `rgba(255, 255, 255, ${flake.opacity})`;
        this.ctx.fill();
        this.ctx.closePath();

        // Thêm hiệu ứng glow
        if (flake.radius > 2) {
            this.ctx.beginPath();
            this.ctx.arc(flake.x, flake.y, flake.radius + 1, 0, Math.PI * 2);
            this.ctx.fillStyle = `rgba(255, 255, 255, ${flake.opacity * 0.3})`;
            this.ctx.fill();
            this.ctx.closePath();
        }
    }

    animate() {
        this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);

        this.snowflakes.forEach(flake => {
            this.updateSnowflake(flake);
            this.drawSnowflake(flake);
        });

        requestAnimationFrame(() => this.animate());
    }
}

// Khởi tạo khi DOM ready
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', () => {
        new InteractiveSnowfall();
    });
} else {
    new InteractiveSnowfall();
}
