using System;
using System.Collections;
using System.Collections.Generic;
using System.Data;
using System.Diagnostics;
using System.IO;
using System.IO.Compression;
using System.Linq;
using System.Runtime.InteropServices;
using System.Runtime.Serialization;
using System.Runtime.Serialization.Json;
using System.Windows.Forms;

namespace Sav
{
    public partial class Form : System.Windows.Forms.Form
    {
        readonly string configPath;
        readonly Config config;
        string selectedGameName;
        Game selectedGame;
        readonly DataContractJsonSerializer serializer = new DataContractJsonSerializer(typeof(Config), new DataContractJsonSerializerSettings()
        {
            UseSimpleDictionaryFormat = true
        });

        public Form()
        {
            InitializeComponent();
            CenterToScreen();

            try
            {
                var dir = Environment.GetFolderPath(Environment.SpecialFolder.ApplicationData);
                dir = Path.Combine(dir, "sav");
                configPath = Path.Combine(dir, "config.json");
                if (File.Exists(configPath))
                {
                    using (var file = File.OpenRead(configPath))
                        config = serializer.ReadObject(file) as Config;
                }
                else
                {
                    Directory.CreateDirectory(dir);
                    config = new Config();
                }
            }
            catch (Exception ex)
            {
                MessageBox.Show(ex.Message);
                Environment.Exit(1);
            }
            if (config.Games == null)
                config.Games = new Dictionary<string, Game>();
            foreach (var game in config.Games.Keys)
                gameList.Items.Add(game);
        }

        private void GameList_SelectedIndexChanged(object sender, EventArgs e)
        {
            var gameSelected = gameList.SelectedItem != null;
            buttonRemoveGame.Enabled = tableLayoutPanelGameManager.Visible = buttonBackUp.Visible = gameSelected;
            if (!gameSelected) return;

            selectedGameName = gameList.SelectedItem.ToString();
            selectedGame = config.Games[selectedGameName];
            textBoxGameName.Text = selectedGameName;
            textBoxSaveRoot.Text = selectedGame.SaveRoot;
            textBoxBackupPath.Text = selectedGame.BackupPath;
            UpdateFileTree();
        }

        private void ButtonAddGame_Click(object sender, EventArgs e)
        {
            for (var i = 1; ; i++)
            {
                var name = $"Game {i}";
                if (!config.Games.ContainsKey(name))
                {
                    try
                    {
                        config.Games[name] = new Game();
                        SaveConfig();
                        gameList.Items.Add(name);
                        gameList.SelectedItem = name;
                    }
                    catch (Exception ex)
                    {
                        config.Games.Remove(name);
                        MessageBox.Show(ex.Message);
                    }
                    break;
                }
            }
        }

        private void ButtonRemoveGame_Click(object sender, EventArgs e)
        {
            if (MessageBox.Show($"Are you sure to delete {selectedGameName}?", "", MessageBoxButtons.YesNo) == DialogResult.Yes)
            {
                try
                {
                    config.Games.Remove(selectedGameName);
                    SaveConfig();
                    gameList.Items.Remove(selectedGameName);
                }
                catch (Exception ex)
                {
                    config.Games[selectedGameName] = selectedGame;
                    MessageBox.Show(ex.Message);
                }
            }
        }

        private void ButtonSaveName_Click(object sender, EventArgs e)
        {
            var newName = textBoxGameName.Text;
            if (config.Games.ContainsKey(newName))
            {
                MessageBox.Show("Game name already exists.");
                return;
            }
            try
            {
                config.Games.Remove(selectedGameName);
                config.Games[newName] = selectedGame;
                SaveConfig();
                gameList.Items[gameList.Items.IndexOf(selectedGameName)] = newName;
            }
            catch (Exception ex)
            {
                config.Games.Remove(newName);
                config.Games[selectedGameName] = selectedGame;
                MessageBox.Show(ex.Message);
            }
        }

        private void ButtonChooseSaveRoot_Click(object sender, EventArgs e)
        {
            using (var dialog = new FolderBrowserDialog())
            {
                dialog.SelectedPath = selectedGame.SaveRoot;
                if (dialog.ShowDialog() == DialogResult.OK)
                {
                    var oldSaveRoot = selectedGame.SaveRoot;
                    try
                    {
                        selectedGame.SaveRoot = dialog.SelectedPath;
                        SaveConfig();
                        textBoxSaveRoot.Text = dialog.SelectedPath;
                        UpdateFileTree();
                    }
                    catch (Exception ex)
                    {
                        selectedGame.SaveRoot = oldSaveRoot;
                        MessageBox.Show(ex.Message);
                    }
                }
            }
        }

        private void ButtonOpenSaveRoot_Click(object sender, EventArgs e)
        {
            if (selectedGame.SaveRoot == null)
            {
                MessageBox.Show("No save root specified.");
            }
            else
            {
                try
                {
                    Process.Start("explorer.exe", selectedGame.SaveRoot);
                }
                catch (Exception ex)
                {
                    MessageBox.Show(ex.Message);
                }
            }
        }

        private void ButtonChooseBackupPath_Click(object sender, EventArgs e)
        {
            using (var dialog = new FolderBrowserDialog())
            {
                dialog.SelectedPath = selectedGame.BackupPath;
                if (dialog.ShowDialog() == DialogResult.OK)
                {
                    var oldBackupPath = selectedGame.BackupPath;
                    try
                    {
                        selectedGame.BackupPath = dialog.SelectedPath;
                        SaveConfig();
                        textBoxBackupPath.Text = dialog.SelectedPath;
                        buttonBackUp.Enabled = selectedGame.SaveRoot != null;
                    }
                    catch (Exception ex)
                    {
                        selectedGame.BackupPath = oldBackupPath;
                        MessageBox.Show(ex.Message);
                    }
                }
            }
        }

        private void ButtonOpenBackup_Click(object sender, EventArgs e)
        {
            if (selectedGame.BackupPath == null)
            {
                MessageBox.Show("No backup path specified.");
            }
            else
            {
                try
                {
                    var fileName = string.Concat(selectedGameName.Split(Path.GetInvalidFileNameChars()));
                    var backup = Path.Combine(selectedGame.BackupPath, $"{fileName}.zip");
                    Process.Start("explorer.exe", File.Exists(backup) ? $"/select,\"{backup}\"" : selectedGame.BackupPath);
                }
                catch (Exception ex)
                {
                    MessageBox.Show(ex.Message);
                }
            }
        }

        private void TreeView_AfterCheck(object sender, TreeViewEventArgs e)
        {
            if (e.Action == TreeViewAction.Unknown) return;

            foreach (var node in GetNodes(e.Node))
                node.Checked = e.Node.Checked;
            if (e.Node.Checked)
            {
                for (var cur = e.Node.Parent; cur != null; cur = cur.Parent)
                    if (cur.Nodes.Cast<TreeNode>().All(n => n.Checked))
                        cur.Checked = true;
            }
            else
            {
                for (var cur = e.Node.Parent; cur != null; cur = cur.Parent)
                    cur.Checked = false;
            }

            var root = treeView.Nodes[0];
            if (root.Checked)
            {
                selectedGame.Paths = new List<string> { "*" };
            }
            else
            {
                IEnumerable<(TreeNode, string)> GetCheckedNodes(TreeNode node, string parent)
                {
                    foreach (TreeNode n in node.Nodes)
                    {
                        var path = Path.Combine(parent, n.Text);
                        if (n.Checked)
                            yield return (n, path);
                        else
                            foreach (var child in GetCheckedNodes(n, path))
                                yield return child;
                    }
                }
                if (selectedGame.Paths == null)
                    selectedGame.Paths = new List<string>();
                else
                    selectedGame.Paths.Clear();
                foreach (var (_, path) in GetCheckedNodes(root, ""))
                    selectedGame.Paths.Add(path);
            }

            try
            {
                SaveConfig();
            }
            catch (Exception ex)
            {
                MessageBox.Show(ex.Message);
            }
        }

        private void ButtonBackUp_Click(object sender, EventArgs e)
        {
            if (selectedGame.Paths == null || selectedGame.Paths.Count == 0)
            {
                MessageBox.Show("No file selected.");
                return;
            }

            var temp = Path.GetTempFileName();
            try
            {
                using (var fs = new FileStream(temp, FileMode.Create))
                using (var archive = new ZipArchive(fs, ZipArchiveMode.Create))
                {
                    if (selectedGame.Paths.Contains("*"))
                        foreach (var file in Directory.EnumerateFiles(selectedGame.SaveRoot, "*", SearchOption.AllDirectories))
                            archive.CreateEntryFromFile(file, file.Substring(selectedGame.SaveRoot.Length + 1));
                    else
                        foreach (var path in selectedGame.Paths)
                        {
                            var name = Path.Combine(selectedGame.SaveRoot, path);
                            if (Directory.Exists(name))
                                foreach (var file in Directory.EnumerateFiles(name, "*", SearchOption.AllDirectories))
                                    archive.CreateEntryFromFile(file, file.Substring(selectedGame.SaveRoot.Length + 1));
                            else
                                archive.CreateEntryFromFile(name, path);
                        }
                }
                var fileName = string.Concat(selectedGameName.Split(Path.GetInvalidFileNameChars()));
                var backup = Path.Combine(selectedGame.BackupPath, $"{fileName}.zip");
                File.Delete(backup);
                File.Move(temp, backup);
                MessageBox.Show("Done.");
            }
            catch (Exception ex)
            {
                File.Delete(temp);
                MessageBox.Show(ex.Message);
            }
        }

        void UpdateFileTree()
        {
            treeView.Nodes.Clear();
            buttonBackUp.Enabled = selectedGame.SaveRoot != null && selectedGame.BackupPath != null;
            if (selectedGame.SaveRoot == null || !Directory.Exists(selectedGame.SaveRoot))
                return;

            var root = treeView.Nodes.Add(new DirectoryInfo(selectedGame.SaveRoot).Name);
            foreach (var file in Directory.EnumerateFiles(selectedGame.SaveRoot, "*", SearchOption.AllDirectories))
            {
                var path = file.Substring(selectedGame.SaveRoot.Length + 1);
                var cur = root;
                foreach (var part in path.Split(new[] { Path.DirectorySeparatorChar, Path.AltDirectorySeparatorChar }))
                {
                    var target = cur.Nodes.Cast<TreeNode>().Where(n => n.Text == part).FirstOrDefault();
                    cur = target ?? cur.Nodes.Add(part);
                }
            }
            treeView.TreeViewNodeSorter = new NodeSorter();
            root.Expand();

            if (selectedGame.Paths == null) return;
            if (selectedGame.Paths.Contains("*"))
            {
                root.Checked = true;
                foreach (var node in GetNodes(root))
                    node.Checked = true;
            }
            else
            {
                foreach (var path in selectedGame.Paths)
                {
                    var cur = root;
                    foreach (var part in path.Split(new[] { Path.DirectorySeparatorChar, Path.AltDirectorySeparatorChar }))
                    {
                        cur = cur.Nodes.Cast<TreeNode>().Where(n => n.Text == part).FirstOrDefault();
                        if (cur == null)
                            break;
                        cur.Expand();
                    }
                    if (cur == null)
                        continue;
                    cur.Checked = true;
                    foreach (var n in GetNodes(cur))
                        n.Checked = true;
                }
            }
        }

        void SaveConfig()
        {
            using (var file = File.Create(configPath))
                serializer.WriteObject(file, config);
        }

        static IEnumerable<TreeNode> GetNodes(TreeNode root)
        {
            foreach (TreeNode node in root.Nodes)
            {
                yield return node;
                foreach (var child in GetNodes(node))
                    yield return child;
            }
        }
    }

    class NodeSorter : IComparer
    {
        [DllImport("shlwapi.dll", CharSet = CharSet.Unicode)]
        static extern int StrCmpLogicalW(string psz1, string psz2);

        public int Compare(object x, object y)
        {
            var tx = x as TreeNode;
            var ty = y as TreeNode;

            if ((tx.Nodes.Count == 0) == (ty.Nodes.Count == 0))
                return StrCmpLogicalW(tx.Text, ty.Text);
            return tx.Nodes.Count == 0 ? 1 : -1;
        }
    }

    [DataContract]
    class Game
    {
        [DataMember]
        public string SaveRoot;
        [DataMember]
        public string BackupPath;
        [DataMember]
        public List<string> Paths;
    }

    [DataContract]
    class Config
    {
        [DataMember]
        public Dictionary<string, Game> Games;
    }
}
