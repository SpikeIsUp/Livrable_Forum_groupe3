let utilisateurConnecte = null;

// ==========================================
// INITIALISATION ET THÈME
// ==========================================
document.addEventListener('DOMContentLoaded', () => {
    changerVue('vue-accueil');
    verifierEtatConnexion();

    // Vérifie la mémoire du navigateur pour le mode sombre
    if (localStorage.getItem('theme') === 'sombre') {
        document.body.classList.add('mode-sombre');
    }
});

function changerTheme() {
    // Bascule la classe sur le body
    document.body.classList.toggle('mode-sombre');
    
    // Sauvegarde le choix de l'utilisateur
    if (document.body.classList.contains('mode-sombre')) {
        localStorage.setItem('theme', 'sombre');
    } else {
        localStorage.setItem('theme', 'clair');
    }
}

function toggleSidebar() {
    document.getElementById('sidebar-menu').classList.toggle('retracte');
    document.getElementById('app-layout').classList.toggle('menu-ferme');
}

// ==========================================
// NAVIGATION ET ÉTAT DE CONNEXION
// ==========================================
function changerVue(idVue) {
    const vues = [
        'vue-accueil', 'vue-nouveau-sujet', 'vue-sujet', 
        'vue-connexion', 'vue-inscription', 'vue-profil',
        'vue-tendance', 'vue-suivis', 'vue-messages','vue-lecture-post'
    ];

    vues.forEach(vue => {
        const element = document.getElementById(vue);
        if (element) element.style.display = 'none';
    });

    const vueActive = document.getElementById(idVue);
    if (vueActive) vueActive.style.display = 'block';

    // Chargements automatiques selon la page
    if (idVue === 'vue-accueil') chargerCategories();
    if (idVue === 'vue-nouveau-sujet') chargerCategoriesDansSelect();
    if (idVue === 'vue-profil') afficherProfil();
    if (idVue === 'vue-messages') chargerBoiteReception();
    if (idVue === 'vue-tendance') chargerTendances();
    if (idVue === 'vue-lecture-post') chargerPostEtCommentaires();
}

function verifierEtatConnexion() {
    const zoneAuth = document.getElementById('zone-auth');
    const msgBtn = document.getElementById('btn-nav-messages');
    const creerBtn = document.getElementById('btn-nav-creer');

    if (utilisateurConnecte) {
        // Affiche un bel avatar généré dynamiquement avec la première lettre du pseudo
        zoneAuth.innerHTML = `<img src="https://ui-avatars.com/api/?name=${utilisateurConnecte.pseudo}&background=1d4ed8&color=fff" class="avatar-mini" onclick="changerVue('vue-profil')" alt="Profil">`;
        if(msgBtn) msgBtn.style.display = 'flex';
        if(creerBtn) creerBtn.style.display = 'flex';
    } else {
        zoneAuth.innerHTML = `<button class="btn-action" onclick="changerVue('vue-connexion')">Connexion</button>`;
        if(msgBtn) msgBtn.style.display = 'none';
        if(creerBtn) creerBtn.style.display = 'none';
    }
}

// ==========================================
// AUTHENTIFICATION (Connexion, Inscription)
// ==========================================
async function tenterConnexion() {
    const pseudo = document.getElementById('login-pseudo').value;
    const mdp = document.getElementById('login-mdp').value;

    if (!pseudo || !mdp) return;

    try {
        const reponse = await fetch('http://localhost:8080/api/connexion', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ pseudo: pseudo, mdp: mdp })
        });

        if (reponse.ok) {
            utilisateurConnecte = await reponse.json();
            document.getElementById('login-pseudo').value = '';
            document.getElementById('login-mdp').value = '';
            verifierEtatConnexion();
            changerVue('vue-accueil');
        } else {
            alert("❌ Pseudo ou mot de passe incorrect.");
        }
    } catch (erreur) {
        alert("Impossible de joindre le serveur.");
    }
}

async function tenterInscription() {
    const pseudo = document.getElementById('insc-pseudo').value;
    const mdp = document.getElementById('insc-mdp').value;

    if (!pseudo || !mdp) {
        alert("Remplis tous les champs !");
        return;
    }

    try {
        const reponse = await fetch('http://localhost:8080/api/inscription', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ pseudo: pseudo, mdp: mdp })
        });

        if (reponse.ok) {
            alert("✅ Compte créé avec succès ! Tu peux maintenant te connecter.");
            document.getElementById('insc-pseudo').value = '';
            document.getElementById('insc-mdp').value = '';
            changerVue('vue-connexion');
        } else if (reponse.status === 409) {
            alert("❌ Ce pseudo est déjà pris, choisis-en un autre.");
        } else {
            alert("❌ Erreur lors de l'inscription.");
        }
    } catch (erreur) {
        alert("Impossible de joindre le serveur.");
    }
}

function seDeconnecter() {
    utilisateurConnecte = null;
    verifierEtatConnexion();
    changerVue('vue-accueil');
}

function afficherProfil() {
    if (!utilisateurConnecte) return;
    document.getElementById('profil-nom').innerText = utilisateurConnecte.pseudo;
}

// ==========================================
// FORUM : CATÉGORIES ET SUJETS
// ==========================================
async function chargerCategories() {
    const conteneur = document.getElementById('liste-categories');
    conteneur.innerHTML = '<p>Chargement...</p>';
    try {
        const reponse = await fetch('http://localhost:8080/api/categories');
        const categories = await reponse.json();
        conteneur.innerHTML = '';
        
        categories.forEach(cat => {
            const div = document.createElement('div');
            div.className = 'carte-forum';
            
            const h2 = document.createElement('h2');
            h2.innerText = cat.titre;
            
            const btn = document.createElement('button');
            btn.className = 'btn-action';
            btn.style.marginTop = '10px';
            btn.innerText = "Voir les discussions";
            btn.onclick = () => chargerSujetsDeCategorie(cat.id, cat.titre);
            
            div.appendChild(h2);
            div.appendChild(btn);
            conteneur.appendChild(div);
        });
    } catch (erreur) {
        conteneur.innerHTML = "<p style='color: red;'>Erreur de connexion au serveur.</p>";
    }
}

async function chargerCategoriesDansSelect() {
    const select = document.getElementById('select-categorie');
    try {
        const reponse = await fetch('http://localhost:8080/api/categories');
        const categories = await reponse.json();
        select.innerHTML = '<option value="">-- Sélectionne une catégorie --</option>';
        categories.forEach(cat => {
            select.innerHTML += `<option value="${cat.id}">${cat.titre}</option>`;
        });
    } catch (erreur) {}
}

async function chargerSujetsDeCategorie(idCategorie, titreCategorie) {
    changerVue('vue-sujet');
    document.getElementById('titre-sujet-actif').innerText = "Sujets : " + titreCategorie;
    const conteneur = document.getElementById('messages-sujet');
    conteneur.innerHTML = '<p>Recherche des posts en cours...</p>';

    try {
        const reponse = await fetch(`http://localhost:8080/api/categories/posts?id=${idCategorie}`);
        const posts = await reponse.json();
        conteneur.innerHTML = '';

        if (!posts || posts.length === 0) {
            conteneur.innerHTML = '<p style="color: var(--text-muted); font-style: italic;">Aucun post dans cette catégorie pour le moment. Sois le premier !</p>';
            return;
        }

        posts.forEach(post => {
            const div = document.createElement('div');
            div.className = 'carte-forum';

            div.innerHTML = `
                <div style="display: flex; justify-content: space-between; align-items: start;">
                    <h3 style="color: var(--primary); margin-bottom: 8px;">${post.titre}</h3>
                    <div style="display: flex; gap: 10px;">
                        <button onclick="reagirPost(${post.id}, 1)" style="background: var(--bg-main); padding: 5px 12px; border-radius: 20px; font-weight: bold; font-size: 14px; border: 1px solid var(--border-color); color: #16a34a; cursor: pointer; transition: 0.2s;">👍 ${post.nb_likes || 0}</button>
                        <button onclick="reagirPost(${post.id}, -1)" style="background: var(--bg-main); padding: 5px 12px; border-radius: 20px; font-weight: bold; font-size: 14px; border: 1px solid var(--border-color); color: #dc2626; cursor: pointer; transition: 0.2s;">👎 ${post.nb_dislikes || 0}</button>
                    </div>
                </div>
                <small style="color: var(--text-muted); display: block; margin-bottom: 15px; border-bottom: 1px solid var(--border-color); padding-bottom: 10px;">
                    Posté par <strong>${post.auteur}</strong> le ${post.date}
                </small>
                <p style="line-height: 1.6; white-space: pre-wrap;">${post.corps}</p>
            `;

            const btn = document.createElement('button');
            btn.className = 'btn-action';
            btn.style.marginTop = '15px';
            btn.innerText = 'Lire et répondre 💬';
            
            btn.onclick = () => ouvrirPost(post.id, post.titre, post.corps, post.auteur, post.date, post.nb_likes, post.nb_dislikes);

            div.appendChild(btn);
            conteneur.appendChild(div);
        });

    } catch (erreur) {
        console.error("Détail de l'erreur :", erreur);
        conteneur.innerHTML = "<p style='color: red;'>Erreur lors du chargement des posts.</p>";
    }
}

async function publierSujet() {
    if (!utilisateurConnecte) {
        alert("Tu dois être connecté pour publier.");
        return;
    }
    const titre = document.getElementById('titre-sujet').value;
    const contenu = document.getElementById('contenu-sujet').value;
    const idCategorie = document.getElementById('select-categorie').value; 
    
    if (!titre || !contenu || !idCategorie) {
        alert("Remplis tous les champs du formulaire.");
        return;
    }

    const donnees = {
        titre: titre,
        corps: contenu,
        id_categorie: parseInt(idCategorie),
        id_utilisateur: utilisateurConnecte.id
    };

try {
        const reponse = await fetch('http://localhost:8080/api/sujets', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(donnees)
        });

        if (reponse.ok) {
            alert("✅ Sujet publié !");
            document.getElementById('titre-sujet').value = '';
            document.getElementById('contenu-sujet').value = '';
            changerVue('vue-accueil');
        } else {
            // NOUVEAU : On affiche une erreur si le serveur refuse l'enregistrement
            alert("❌ Le serveur Go a refusé le post. Il y a un problème côté base de données.");
        }
    } catch (erreur) {
        // NOUVEAU : On affiche une erreur si le serveur est éteint
        alert("❌ Impossible de joindre le serveur Go. Est-il bien allumé ?");
        console.error(erreur);
    }

// ==========================================
// TENDANCES ET LIKES
// ==========================================
async function chargerTendances() {
    const box = document.getElementById('liste-tendances');
    box.innerHTML = '<p>Chargement des tendances...</p>';
    try {
        const res = await fetch('http://localhost:8080/api/tendances');
        const posts = await res.json();
        box.innerHTML = '';
        
        if(!posts || posts.length === 0) {
            box.innerHTML = '<p>Aucune tendance pour le moment.</p>';
            return;
        }

        posts.forEach(p => {
            box.innerHTML += `
                <div class="carte-forum">
                    <h3 style="color:var(--primary)">${p.titre}</h3>
                    <small style="color:var(--text-muted)">Par ${p.auteur}</small>
                    <p style="margin-top:10px;">${p.corps}</p>
                    <button class="btn-action" style="margin-top:15px; background:var(--accent)" onclick="likerPost(${p.id})">
                        ❤️ ${p.nb_likes} Likes
                    </button>
                </div>`;
        });
    } catch(e) {
        box.innerHTML = "<p style='color: red;'>Erreur serveur.</p>";
    }
}

async function likerPost(idPost) {
    if(!utilisateurConnecte) { 
        alert("Connecte-toi pour liker des posts !"); 
        return; 
    }
    try {
        await fetch(`http://localhost:8080/api/like?uid=${utilisateurConnecte.id}&pid=${idPost}`);
        chargerTendances(); // Rafraîchit les likes à l'écran
    } catch(e) {}
}

// ==========================================
// MESSAGERIE PRIVÉE
// ==========================================
async function envoyerMP() {
    if (!utilisateurConnecte) return;

    const destinataire = document.getElementById('mp-destinataire').value;
    const contenu = document.getElementById('mp-contenu').value;

    if (!destinataire || !contenu) {
        alert("Remplis le pseudo et le message !");
        return;
    }

    const donnees = {
        id_expediteur: utilisateurConnecte.id,
        pseudo_destinataire: destinataire,
        contenu: contenu
    };

    try {
        const reponse = await fetch('http://localhost:8080/api/messages/envoyer', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(donnees)
        });

        if (reponse.ok) {
            alert("✅ Message envoyé avec succès !");
            document.getElementById('mp-destinataire').value = '';
            document.getElementById('mp-contenu').value = '';
            chargerBoiteReception(); // Recharge la boîte après l'envoi
        } else if (reponse.status === 404) {
            alert("❌ Cet utilisateur n'existe pas. Vérifie le pseudo.");
        } else {
            alert("❌ Erreur lors de l'envoi.");
        }
    } catch (erreur) {
        console.error(erreur);
    }
}

async function chargerBoiteReception() {
    if (!utilisateurConnecte) return;

    const conteneur = document.getElementById('liste-messages');
    conteneur.innerHTML = '<p>Recherche des messages...</p>';

    try {
        const reponse = await fetch(`http://localhost:8080/api/messages/lire?id=${utilisateurConnecte.id}`);
        const messages = await reponse.json();

        conteneur.innerHTML = '';
        if (!messages || messages.length === 0) {
            conteneur.innerHTML = '<p style="color: var(--text-muted); font-style: italic;">Ta boîte de réception est vide.</p>';
            return;
        }

        messages.forEach(msg => {
            const div = document.createElement('div');
            div.className = 'carte-forum';
            div.style.borderLeft = '4px solid var(--primary)';

            div.innerHTML = `
                <div style="display: flex; justify-content: space-between; margin-bottom: 10px;">
                    <strong>De : ${msg.expediteur}</strong>
                    <small style="color: var(--text-muted);">${msg.date_envoi}</small>
                </div>
                <p>${msg.contenu}</p>
            `;
            conteneur.appendChild(div);
        });
    } catch (erreur) {
        conteneur.innerHTML = "<p style='color: red;'>Impossible de charger les messages.</p>";
    }
}  
}
// ==========================================
// LECTURE D'UN POST ET COMMENTAIRES
// ==========================================
let postActifId = null;

function ouvrirPost(id, titre, corps, auteur, date, nbLikes, nbDislikes) {
    changerVue('vue-lecture-post');
    postActifId = id;
    
    document.getElementById('post-original-titre').innerText = titre;
    document.getElementById('post-original-auteur').innerHTML = `Posté par <strong>${auteur}</strong> le ${date}`;
    document.getElementById('post-original-corps').innerText = corps;
    
    // On injecte les boutons au bon endroit
    const zoneBoutons = document.getElementById('zone-reactions-post');
    zoneBoutons.innerHTML = `
        <div style="display: flex; gap: 15px; margin-top: 20px;">
            <button id="btn-like-post" class="btn-action" style="background: #16a34a; font-weight: bold; border-radius: 20px; width: 120px;">
                👍 ${nbLikes || 0}
            </button>
            <button id="btn-dislike-post" class="btn-action" style="background: #dc2626; font-weight: bold; border-radius: 20px; width: 120px;">
                👎 ${nbDislikes || 0}
            </button>
        </div>
    `;

    // On active les clics
    document.getElementById('btn-like-post').onclick = () => reagirPost(id, 1);
    document.getElementById('btn-dislike-post').onclick = () => reagirPost(id, -1);

    chargerCommentaires();
}
async function reagirPost(idPost, typeReaction) {
    if (!utilisateurConnecte) {
        alert("❌ Tu dois être connecté pour liker ou disliker !");
        return;
    }
    
    try {
        // NOUVEAU : On utilise la nouvelle adresse "anti-adblock"
        const url = `http://localhost:8080/api/reaction-post?uid=${utilisateurConnecte.id}&pid=${idPost}&type=${typeReaction}`;
        const reponse = await fetch(url);
        
        if (reponse.ok) {
            const btnLike = document.getElementById('btn-like-post');
            const btnDislike = document.getElementById('btn-dislike-post');
            
            if (btnLike && typeReaction === 1) btnLike.style.background = '#15803d';
            if (btnDislike && typeReaction === -1) btnDislike.style.background = '#b91c1c';
            
            changerVue('vue-accueil'); 
        } else {
            alert("❌ Le serveur Go a refusé le vote. Code : " + reponse.status);
        }
    } catch (erreur) { 
        // NOUVEAU : On affiche l'erreur technique exacte du navigateur
        alert("❌ Erreur réseau interceptée par le navigateur : " + erreur.message);
        console.error(erreur);
    }
}

async function chargerCommentaires() {
    const box = document.getElementById('liste-commentaires');
    box.innerHTML = '<p>Chargement des réponses...</p>';
    
    try {
        const res = await fetch(`http://localhost:8080/api/commentaires/lire?id_post=${postActifId}`);
        const comms = await res.json();
        
        box.innerHTML = '';
        if(!comms || comms.length === 0) {
            box.innerHTML = '<p style="color: var(--text-muted); font-style: italic;">Aucune réponse pour le moment. Sois le premier !</p>';
            return;
        }
        
        comms.forEach(c => {
            box.innerHTML += `
                <div class="carte-forum" style="padding: 15px; margin-bottom: 10px; border-left: 3px solid var(--accent);">
                    <div style="display: flex; justify-content: space-between; margin-bottom: 5px;">
                        <strong>${c.auteur}</strong>
                        <small style="color: var(--text-muted);">${c.date}</small>
                    </div>
                    <p style="margin-top: 5px;">${c.contenu}</p>
                </div>`;
        });
    } catch(e) {
        box.innerHTML = '<p style="color: red;">Erreur de chargement des commentaires.</p>';
    }
}

async function envoyerCommentaire() {
    if(!utilisateurConnecte) {
        alert("Tu dois être connecté pour répondre.");
        return;
    }
    
    const txt = document.getElementById('nouveau-commentaire-texte').value;
    if(!txt) return;

    try {
        const reponse = await fetch('http://localhost:8080/api/commentaires', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({
                id_post: postActifId, 
                id_utilisateur: utilisateurConnecte.id, 
                contenu: txt
            })
        });
        
        if(reponse.ok) {
            document.getElementById('nouveau-commentaire-texte').value = '';
            chargerCommentaires(); // On rafraîchit la liste pour voir sa réponse !
        }
    } catch(e) {
        alert("Erreur lors de l'envoi.");
    }
}

async function envoyerCommentaire() {
    // 1. Vérification de la connexion
    if(!utilisateurConnecte) {
        alert("Tu dois être connecté pour répondre.");
        return;
    }
    
    const txt = document.getElementById('nouveau-commentaire-texte').value;
    if(!txt) {
        alert("Tu ne peux pas envoyer un message vide !");
        return;
    }

    try {
        const reponse = await fetch('http://localhost:8080/api/commentaires', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({
                id_post: parseInt(postActifId), // On s'assure que c'est bien un chiffre
                id_utilisateur: utilisateurConnecte.id, 
                contenu: txt
            })
        });
        
        if(reponse.ok) {
            document.getElementById('nouveau-commentaire-texte').value = '';
            chargerCommentaires(); 
        } else {
            // NOUVEAU : On capture l'erreur côté base de données
            alert("❌ Le serveur a refusé d'enregistrer la réponse. As-tu bien créé la table 'commentaire' dans phpMyAdmin ?");
        }
    } catch(e) {
        // NOUVEAU : On capture l'erreur de serveur éteint
        alert("❌ Impossible de joindre le serveur Go.");
        console.error(e);
    }
}