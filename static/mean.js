var meanout = function(elem){
    var cs = getComputedStyle(elem);

    var c = elem.getContext('2d');
    c.strokeStyle = 'rgb(40,40,40)';
    c.lineWidth = 2;
    c.translate(0.5,0.5)
    c.lineCap = 'round';

    var currentMean = elem.dataset.val;
    var draw = function(){
        var w = parseInt(cs.getPropertyValue('width'), 10);
        var h = parseInt(cs.getPropertyValue('height'), 10);
        elem.height = h;
        elem.width = w;

        // Clear
        c.save();
        c.fillStyle = 'rgb(255,255,255)';
        c.rect(0, 0, w, h);
        c.fill();
        c.closePath();
        c.restore();

        // Line
        var lineStart = {x: 1, y: h/2}
        var lineEnd = {x: w-1, y: h/2};
        c.save();
        c.beginPath();
        c.moveTo(lineStart.x, lineStart.y);
        c.lineTo(lineEnd.x, lineEnd.y);
        c.stroke();
        c.closePath();
        c.restore();

        // Left cap
        c.save();
        c.beginPath();
        c.moveTo(lineStart.x, lineStart.y-(h/4));
        c.lineTo(lineStart.x, lineStart.y+(h/4))
        c.stroke();
        c.closePath();
        c.restore();

        // Right cap
        c.save();
        c.beginPath();
        c.moveTo(lineEnd.x, lineEnd.y-(h/4));
        c.lineTo(lineEnd.x, lineEnd.y+(h/4))
        c.stroke();
        c.closePath();
        c.restore();

        // Pointer
        var targetMean = elem.dataset.val;
        if (currentMean < targetMean){
            currentMean += (targetMean - currentMean) / 20;
        } else {
            currentMean -= (currentMean - targetMean) / 20;
        }
        var pointerX = (currentMean/1000000000)*w;

        c.save();
        c.beginPath();
        c.moveTo(pointerX, h/2-3);
        c.lineTo(pointerX-7, h/2-15)
        c.lineTo(pointerX+7, h/2-15)
        c.lineTo(pointerX, h/2-3);
        c.strokeStyle = 'rgb(220, 40, 40)';
        c.stroke();
        c.fillStyle = 'rgb(220, 40, 40)';
        c.fill();
        c.closePath();
        c.restore();

        requestAnimationFrame(draw);
    };
    draw();
};
