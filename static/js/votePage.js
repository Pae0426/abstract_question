$('.vote-page-modal-btn').on('click', function() {
    $('.vote-page-modal-item').fadeIn();
    $.ajax({
        dataType: 'json',
        contentType: 'application/json',
        type: 'GET',
        url: '/get-vote-page-info',
    }).done(function(vote_info) {
        //ページ投票グラフ表示
        let data_label = [];
        let data_color = [];
        for(let i=0;i<=page_total;i++) {
            if(vote_info.userVotePage[i] == 1) {
                data_color[i] = '#9eff9e';
            } else {
                data_color[i] = '#9eceff';
            }
        }
        for(let i=1;i<=page_total;i++) {
            data_label.push(i.toString());
        }
        let ctx = document.getElementById("chart").getContext('2d');
        myChart = new Chart(ctx, {
            type: "bar",
            data: {
                labels: data_label,
                datasets: [
                    {
                        label: "投票数",
                        data: vote_info.voteCount,
                        backgroundColor: data_color,
                    },
                ]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    legend: {
                        display: false
                    }
                },
                scales: {
                    xAxes: [{
                        ticks: {
                            stepSize: 1,
                            autoSkip: false,
                            maxRotation: 0,
                            minRotation: 0,
                        },
                    }],
                    yAxes: [{
                        display: false,
                        ticks: {
                            min: 0,
                            max: 100,
                        },
                    }],
                },
                hover: {
                    mode: 'single'
                },
            }
        });

        let page_now = $('.page-now-text').html();
        console.log(vote_info.userVotePage[page_now-1]);
        if(vote_info.userVotePage[page_now-1] == 1) {
            $('.vote-page-container').hide();
            $('.remove-vote-page-container').show();
        } else {
            $('.remove-vote-page-container').hide();
            $('.vote-page-container').show();
        }
    }).fail(function() {
        console.log('通信失敗');
    });
});

$('.vote-page-modal-close-btn, .vote-page-modal-bg').on('click', function() {
    if (myChart) {
        myChart.destroy();
    }
    $('.vote-page-modal-item').fadeOut();
});

$('.vote-page-btn').on('click', function() {
    let page_now = $('.page-now-text').html();
    page_now = parseInt(page_now);
    $.ajax({
        dataType: 'json',
        contentType: 'application/json',
        type: 'POST',
        url: '/vote-page',
        data: JSON.stringify({
            page: page_now,
        })
    }).done(function() {
        $('.vote-page-container').hide();
        $('.remove-vote-page-container').show();
        myChart.data.datasets[0].data[page_now-1]++;
        myChart.data.datasets[0].backgroundColor[page_now-1] = '#9eff9e';
        myChart.update();
    }).fail(function() {
        console.log('通信失敗');
    });
});

$('.remove-vote-page-btn').on('click', function() {
    let page_now = $('.page-now-text').html();
    $.ajax({
        dataType: 'json',
        contentType: 'application/json',
        type: 'POST',
        url: '/remove-vote-page',
        data: JSON.stringify({
            page: parseInt(page_now),
        })
    }).done(function() {
        $('.remove-vote-page-container').hide();
        $('.vote-page-container').show();
        myChart.data.datasets[0].data[page_now-1]--;
        myChart.data.datasets[0].backgroundColor[page_now-1] = '#9eceff';
        myChart.update();
    }).fail(function() {
        console.log('通信失敗');
    });
});
