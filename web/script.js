const titles = {
  chat: 'Чат',
  homework: 'Домашнее задание',
  quizzes: 'Квизы',
  stats: 'Статистика',
  leaderboard: 'Лидерборд'
};

document.querySelectorAll('.nav-item').forEach(item => {
  item.addEventListener('click', () => {
    document.querySelectorAll('.nav-item').forEach(n => n.classList.remove('active'));
    document.querySelectorAll('.tab-content').forEach(t => t.classList.remove('active'));

    item.classList.add('active');

    const tab = item.dataset.tab;
    document.getElementById('tab-' + tab).classList.add('active');
    document.getElementById('pageTitle').textContent = titles[tab];
  });
});