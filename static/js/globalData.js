const PDF_DIR = '12';
let winWidth = Math.round($(window).width());
let winHeight = Math.round($(window).height());
let page_total = $('.progressbar').data('total-progress');
const PAGE_TOTAL = parseInt(page_total);
let isVisible = true;
let pageChart;
let isHighlight = false;
let pageVoted;
let pageVotedInfo;
let highlightWidth = {};
