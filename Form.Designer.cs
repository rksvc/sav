namespace Sav
{
    partial class Form
    {
        private System.ComponentModel.IContainer components = null;

        protected override void Dispose(bool disposing)
        {
            if (disposing && (components != null))
            {
                components.Dispose();
            }
            base.Dispose(disposing);
        }


        private void InitializeComponent()
        {
            this.labelGames = new System.Windows.Forms.Label();
            this.gameList = new System.Windows.Forms.ListBox();
            this.buttonAddGame = new System.Windows.Forms.Button();
            this.tableLayoutPanelRoot = new System.Windows.Forms.TableLayoutPanel();
            this.tableLayoutPanelGameManager = new System.Windows.Forms.TableLayoutPanel();
            this.labelBackupPath = new System.Windows.Forms.Label();
            this.textBoxBackupPath = new System.Windows.Forms.TextBox();
            this.buttonChooseBackupPath = new System.Windows.Forms.Button();
            this.buttonOpenBackup = new System.Windows.Forms.Button();
            this.labelGameName = new System.Windows.Forms.Label();
            this.textBoxGameName = new System.Windows.Forms.TextBox();
            this.buttonSaveName = new System.Windows.Forms.Button();
            this.labelSaveRoot = new System.Windows.Forms.Label();
            this.textBoxSaveRoot = new System.Windows.Forms.TextBox();
            this.buttonChooseSaveRoot = new System.Windows.Forms.Button();
            this.buttonOpenSaveRoot = new System.Windows.Forms.Button();
            this.treeView = new System.Windows.Forms.TreeView();
            this.tableLayoutPanelGameControl = new System.Windows.Forms.TableLayoutPanel();
            this.buttonRemoveGame = new System.Windows.Forms.Button();
            this.buttonBackUp = new System.Windows.Forms.Button();
            this.tableLayoutPanelRoot.SuspendLayout();
            this.tableLayoutPanelGameManager.SuspendLayout();
            this.tableLayoutPanelGameControl.SuspendLayout();
            this.SuspendLayout();
            // 
            // labelGames
            // 
            this.labelGames.Location = new System.Drawing.Point(3, 0);
            this.labelGames.Name = "labelGames";
            this.labelGames.Size = new System.Drawing.Size(47, 15);
            this.labelGames.TabIndex = 1;
            this.labelGames.Text = "Games";
            // 
            // gameList
            // 
            this.gameList.Anchor = ((System.Windows.Forms.AnchorStyles)((((System.Windows.Forms.AnchorStyles.Top | System.Windows.Forms.AnchorStyles.Bottom) 
            | System.Windows.Forms.AnchorStyles.Left) 
            | System.Windows.Forms.AnchorStyles.Right)));
            this.gameList.FormattingEnabled = true;
            this.gameList.ItemHeight = 15;
            this.gameList.Location = new System.Drawing.Point(3, 18);
            this.gameList.Name = "gameList";
            this.gameList.Size = new System.Drawing.Size(251, 409);
            this.gameList.TabIndex = 2;
            this.gameList.SelectedIndexChanged += new System.EventHandler(this.GameList_SelectedIndexChanged);
            // 
            // buttonAddGame
            // 
            this.buttonAddGame.Anchor = ((System.Windows.Forms.AnchorStyles)((((System.Windows.Forms.AnchorStyles.Top | System.Windows.Forms.AnchorStyles.Bottom) 
            | System.Windows.Forms.AnchorStyles.Left) 
            | System.Windows.Forms.AnchorStyles.Right)));
            this.buttonAddGame.Location = new System.Drawing.Point(3, 3);
            this.buttonAddGame.Name = "buttonAddGame";
            this.buttonAddGame.Size = new System.Drawing.Size(119, 29);
            this.buttonAddGame.TabIndex = 3;
            this.buttonAddGame.Text = "Add";
            this.buttonAddGame.UseVisualStyleBackColor = true;
            this.buttonAddGame.Click += new System.EventHandler(this.ButtonAddGame_Click);
            // 
            // tableLayoutPanelRoot
            // 
            this.tableLayoutPanelRoot.Anchor = ((System.Windows.Forms.AnchorStyles)((((System.Windows.Forms.AnchorStyles.Top | System.Windows.Forms.AnchorStyles.Bottom) 
            | System.Windows.Forms.AnchorStyles.Left) 
            | System.Windows.Forms.AnchorStyles.Right)));
            this.tableLayoutPanelRoot.ColumnCount = 2;
            this.tableLayoutPanelRoot.ColumnStyles.Add(new System.Windows.Forms.ColumnStyle(System.Windows.Forms.SizeType.Percent, 30F));
            this.tableLayoutPanelRoot.ColumnStyles.Add(new System.Windows.Forms.ColumnStyle(System.Windows.Forms.SizeType.Percent, 70F));
            this.tableLayoutPanelRoot.Controls.Add(this.gameList, 0, 1);
            this.tableLayoutPanelRoot.Controls.Add(this.labelGames, 0, 0);
            this.tableLayoutPanelRoot.Controls.Add(this.tableLayoutPanelGameManager, 1, 1);
            this.tableLayoutPanelRoot.Controls.Add(this.tableLayoutPanelGameControl, 0, 2);
            this.tableLayoutPanelRoot.Controls.Add(this.buttonBackUp, 1, 2);
            this.tableLayoutPanelRoot.Location = new System.Drawing.Point(12, 12);
            this.tableLayoutPanelRoot.Name = "tableLayoutPanelRoot";
            this.tableLayoutPanelRoot.RowCount = 3;
            this.tableLayoutPanelRoot.RowStyles.Add(new System.Windows.Forms.RowStyle());
            this.tableLayoutPanelRoot.RowStyles.Add(new System.Windows.Forms.RowStyle(System.Windows.Forms.SizeType.Percent, 100F));
            this.tableLayoutPanelRoot.RowStyles.Add(new System.Windows.Forms.RowStyle());
            this.tableLayoutPanelRoot.Size = new System.Drawing.Size(858, 474);
            this.tableLayoutPanelRoot.TabIndex = 0;
            // 
            // tableLayoutPanelGameManager
            // 
            this.tableLayoutPanelGameManager.Anchor = ((System.Windows.Forms.AnchorStyles)((((System.Windows.Forms.AnchorStyles.Top | System.Windows.Forms.AnchorStyles.Bottom) 
            | System.Windows.Forms.AnchorStyles.Left) 
            | System.Windows.Forms.AnchorStyles.Right)));
            this.tableLayoutPanelGameManager.ColumnCount = 4;
            this.tableLayoutPanelGameManager.ColumnStyles.Add(new System.Windows.Forms.ColumnStyle());
            this.tableLayoutPanelGameManager.ColumnStyles.Add(new System.Windows.Forms.ColumnStyle(System.Windows.Forms.SizeType.Percent, 100F));
            this.tableLayoutPanelGameManager.ColumnStyles.Add(new System.Windows.Forms.ColumnStyle());
            this.tableLayoutPanelGameManager.ColumnStyles.Add(new System.Windows.Forms.ColumnStyle());
            this.tableLayoutPanelGameManager.Controls.Add(this.labelBackupPath, 0, 2);
            this.tableLayoutPanelGameManager.Controls.Add(this.textBoxBackupPath, 1, 2);
            this.tableLayoutPanelGameManager.Controls.Add(this.buttonChooseBackupPath, 2, 2);
            this.tableLayoutPanelGameManager.Controls.Add(this.buttonOpenBackup, 3, 2);
            this.tableLayoutPanelGameManager.Controls.Add(this.labelGameName, 0, 0);
            this.tableLayoutPanelGameManager.Controls.Add(this.textBoxGameName, 1, 0);
            this.tableLayoutPanelGameManager.Controls.Add(this.buttonSaveName, 2, 0);
            this.tableLayoutPanelGameManager.Controls.Add(this.labelSaveRoot, 0, 1);
            this.tableLayoutPanelGameManager.Controls.Add(this.textBoxSaveRoot, 1, 1);
            this.tableLayoutPanelGameManager.Controls.Add(this.buttonChooseSaveRoot, 2, 1);
            this.tableLayoutPanelGameManager.Controls.Add(this.buttonOpenSaveRoot, 3, 1);
            this.tableLayoutPanelGameManager.Controls.Add(this.treeView, 0, 3);
            this.tableLayoutPanelGameManager.Location = new System.Drawing.Point(260, 18);
            this.tableLayoutPanelGameManager.Name = "tableLayoutPanelGameManager";
            this.tableLayoutPanelGameManager.RowCount = 4;
            this.tableLayoutPanelGameManager.RowStyles.Add(new System.Windows.Forms.RowStyle(System.Windows.Forms.SizeType.Absolute, 35F));
            this.tableLayoutPanelGameManager.RowStyles.Add(new System.Windows.Forms.RowStyle(System.Windows.Forms.SizeType.Absolute, 35F));
            this.tableLayoutPanelGameManager.RowStyles.Add(new System.Windows.Forms.RowStyle(System.Windows.Forms.SizeType.Absolute, 35F));
            this.tableLayoutPanelGameManager.RowStyles.Add(new System.Windows.Forms.RowStyle(System.Windows.Forms.SizeType.Percent, 100F));
            this.tableLayoutPanelGameManager.Size = new System.Drawing.Size(595, 412);
            this.tableLayoutPanelGameManager.TabIndex = 0;
            this.tableLayoutPanelGameManager.Visible = false;
            // 
            // labelBackupPath
            // 
            this.labelBackupPath.Anchor = ((System.Windows.Forms.AnchorStyles)((System.Windows.Forms.AnchorStyles.Left | System.Windows.Forms.AnchorStyles.Right)));
            this.labelBackupPath.AutoSize = true;
            this.labelBackupPath.Location = new System.Drawing.Point(3, 80);
            this.labelBackupPath.Name = "labelBackupPath";
            this.labelBackupPath.Size = new System.Drawing.Size(103, 15);
            this.labelBackupPath.TabIndex = 12;
            this.labelBackupPath.Text = "Backup Path:";
            this.labelBackupPath.TextAlign = System.Drawing.ContentAlignment.MiddleCenter;
            // 
            // textBoxBackupPath
            // 
            this.textBoxBackupPath.Anchor = ((System.Windows.Forms.AnchorStyles)((System.Windows.Forms.AnchorStyles.Left | System.Windows.Forms.AnchorStyles.Right)));
            this.textBoxBackupPath.Location = new System.Drawing.Point(112, 75);
            this.textBoxBackupPath.Name = "textBoxBackupPath";
            this.textBoxBackupPath.ReadOnly = true;
            this.textBoxBackupPath.Size = new System.Drawing.Size(302, 25);
            this.textBoxBackupPath.TabIndex = 13;
            // 
            // buttonChooseBackupPath
            // 
            this.buttonChooseBackupPath.Anchor = ((System.Windows.Forms.AnchorStyles)((((System.Windows.Forms.AnchorStyles.Top | System.Windows.Forms.AnchorStyles.Bottom) 
            | System.Windows.Forms.AnchorStyles.Left) 
            | System.Windows.Forms.AnchorStyles.Right)));
            this.buttonChooseBackupPath.Location = new System.Drawing.Point(420, 73);
            this.buttonChooseBackupPath.Name = "buttonChooseBackupPath";
            this.buttonChooseBackupPath.Size = new System.Drawing.Size(91, 29);
            this.buttonChooseBackupPath.TabIndex = 14;
            this.buttonChooseBackupPath.Text = "Choose";
            this.buttonChooseBackupPath.UseVisualStyleBackColor = true;
            this.buttonChooseBackupPath.Click += new System.EventHandler(this.ButtonChooseBackupPath_Click);
            // 
            // buttonOpenBackup
            // 
            this.buttonOpenBackup.Anchor = ((System.Windows.Forms.AnchorStyles)((((System.Windows.Forms.AnchorStyles.Top | System.Windows.Forms.AnchorStyles.Bottom) 
            | System.Windows.Forms.AnchorStyles.Left) 
            | System.Windows.Forms.AnchorStyles.Right)));
            this.buttonOpenBackup.Location = new System.Drawing.Point(517, 73);
            this.buttonOpenBackup.Name = "buttonOpenBackup";
            this.buttonOpenBackup.Size = new System.Drawing.Size(75, 29);
            this.buttonOpenBackup.TabIndex = 15;
            this.buttonOpenBackup.Text = "Open";
            this.buttonOpenBackup.UseVisualStyleBackColor = true;
            this.buttonOpenBackup.Click += new System.EventHandler(this.ButtonOpenBackup_Click);
            // 
            // labelGameName
            // 
            this.labelGameName.Anchor = ((System.Windows.Forms.AnchorStyles)((System.Windows.Forms.AnchorStyles.Left | System.Windows.Forms.AnchorStyles.Right)));
            this.labelGameName.AutoSize = true;
            this.labelGameName.Location = new System.Drawing.Point(3, 10);
            this.labelGameName.Name = "labelGameName";
            this.labelGameName.Size = new System.Drawing.Size(103, 15);
            this.labelGameName.TabIndex = 5;
            this.labelGameName.Text = "Name:";
            // 
            // textBoxGameName
            // 
            this.textBoxGameName.Anchor = ((System.Windows.Forms.AnchorStyles)((System.Windows.Forms.AnchorStyles.Left | System.Windows.Forms.AnchorStyles.Right)));
            this.textBoxGameName.Location = new System.Drawing.Point(112, 5);
            this.textBoxGameName.Name = "textBoxGameName";
            this.textBoxGameName.Size = new System.Drawing.Size(302, 25);
            this.textBoxGameName.TabIndex = 6;
            // 
            // buttonSaveName
            // 
            this.buttonSaveName.Anchor = ((System.Windows.Forms.AnchorStyles)((((System.Windows.Forms.AnchorStyles.Top | System.Windows.Forms.AnchorStyles.Bottom) 
            | System.Windows.Forms.AnchorStyles.Left) 
            | System.Windows.Forms.AnchorStyles.Right)));
            this.tableLayoutPanelGameManager.SetColumnSpan(this.buttonSaveName, 2);
            this.buttonSaveName.Location = new System.Drawing.Point(420, 3);
            this.buttonSaveName.Name = "buttonSaveName";
            this.buttonSaveName.Size = new System.Drawing.Size(172, 29);
            this.buttonSaveName.TabIndex = 7;
            this.buttonSaveName.Text = "Save";
            this.buttonSaveName.UseVisualStyleBackColor = true;
            this.buttonSaveName.Click += new System.EventHandler(this.ButtonSaveName_Click);
            // 
            // labelSaveRoot
            // 
            this.labelSaveRoot.Anchor = ((System.Windows.Forms.AnchorStyles)((System.Windows.Forms.AnchorStyles.Left | System.Windows.Forms.AnchorStyles.Right)));
            this.labelSaveRoot.AutoSize = true;
            this.labelSaveRoot.Location = new System.Drawing.Point(3, 45);
            this.labelSaveRoot.Name = "labelSaveRoot";
            this.labelSaveRoot.Size = new System.Drawing.Size(103, 15);
            this.labelSaveRoot.TabIndex = 8;
            this.labelSaveRoot.Text = "Save Root:";
            // 
            // textBoxSaveRoot
            // 
            this.textBoxSaveRoot.Anchor = ((System.Windows.Forms.AnchorStyles)((System.Windows.Forms.AnchorStyles.Left | System.Windows.Forms.AnchorStyles.Right)));
            this.textBoxSaveRoot.Location = new System.Drawing.Point(112, 40);
            this.textBoxSaveRoot.Name = "textBoxSaveRoot";
            this.textBoxSaveRoot.ReadOnly = true;
            this.textBoxSaveRoot.Size = new System.Drawing.Size(302, 25);
            this.textBoxSaveRoot.TabIndex = 9;
            // 
            // buttonChooseSaveRoot
            // 
            this.buttonChooseSaveRoot.Anchor = ((System.Windows.Forms.AnchorStyles)((((System.Windows.Forms.AnchorStyles.Top | System.Windows.Forms.AnchorStyles.Bottom) 
            | System.Windows.Forms.AnchorStyles.Left) 
            | System.Windows.Forms.AnchorStyles.Right)));
            this.buttonChooseSaveRoot.Location = new System.Drawing.Point(420, 38);
            this.buttonChooseSaveRoot.Name = "buttonChooseSaveRoot";
            this.buttonChooseSaveRoot.Size = new System.Drawing.Size(91, 29);
            this.buttonChooseSaveRoot.TabIndex = 10;
            this.buttonChooseSaveRoot.Text = "Choose";
            this.buttonChooseSaveRoot.UseVisualStyleBackColor = true;
            this.buttonChooseSaveRoot.Click += new System.EventHandler(this.ButtonChooseSaveRoot_Click);
            // 
            // buttonOpenSaveRoot
            // 
            this.buttonOpenSaveRoot.Anchor = ((System.Windows.Forms.AnchorStyles)((((System.Windows.Forms.AnchorStyles.Top | System.Windows.Forms.AnchorStyles.Bottom) 
            | System.Windows.Forms.AnchorStyles.Left) 
            | System.Windows.Forms.AnchorStyles.Right)));
            this.buttonOpenSaveRoot.Location = new System.Drawing.Point(517, 38);
            this.buttonOpenSaveRoot.Name = "buttonOpenSaveRoot";
            this.buttonOpenSaveRoot.Size = new System.Drawing.Size(75, 29);
            this.buttonOpenSaveRoot.TabIndex = 11;
            this.buttonOpenSaveRoot.Text = "Open";
            this.buttonOpenSaveRoot.UseVisualStyleBackColor = true;
            this.buttonOpenSaveRoot.Click += new System.EventHandler(this.ButtonOpenSaveRoot_Click);
            // 
            // treeView
            // 
            this.treeView.Anchor = ((System.Windows.Forms.AnchorStyles)((((System.Windows.Forms.AnchorStyles.Top | System.Windows.Forms.AnchorStyles.Bottom) 
            | System.Windows.Forms.AnchorStyles.Left) 
            | System.Windows.Forms.AnchorStyles.Right)));
            this.treeView.CheckBoxes = true;
            this.tableLayoutPanelGameManager.SetColumnSpan(this.treeView, 4);
            this.treeView.Location = new System.Drawing.Point(3, 108);
            this.treeView.Name = "treeView";
            this.treeView.Size = new System.Drawing.Size(589, 301);
            this.treeView.TabIndex = 16;
            this.treeView.AfterCheck += new System.Windows.Forms.TreeViewEventHandler(this.TreeView_AfterCheck);
            // 
            // tableLayoutPanelGameControl
            // 
            this.tableLayoutPanelGameControl.Anchor = ((System.Windows.Forms.AnchorStyles)((((System.Windows.Forms.AnchorStyles.Top | System.Windows.Forms.AnchorStyles.Bottom) 
            | System.Windows.Forms.AnchorStyles.Left) 
            | System.Windows.Forms.AnchorStyles.Right)));
            this.tableLayoutPanelGameControl.ColumnCount = 2;
            this.tableLayoutPanelGameControl.ColumnStyles.Add(new System.Windows.Forms.ColumnStyle(System.Windows.Forms.SizeType.Percent, 50F));
            this.tableLayoutPanelGameControl.ColumnStyles.Add(new System.Windows.Forms.ColumnStyle(System.Windows.Forms.SizeType.Percent, 50F));
            this.tableLayoutPanelGameControl.Controls.Add(this.buttonAddGame, 0, 0);
            this.tableLayoutPanelGameControl.Controls.Add(this.buttonRemoveGame, 1, 0);
            this.tableLayoutPanelGameControl.Location = new System.Drawing.Point(3, 436);
            this.tableLayoutPanelGameControl.Name = "tableLayoutPanelGameControl";
            this.tableLayoutPanelGameControl.RowCount = 1;
            this.tableLayoutPanelGameControl.RowStyles.Add(new System.Windows.Forms.RowStyle(System.Windows.Forms.SizeType.Percent, 100F));
            this.tableLayoutPanelGameControl.Size = new System.Drawing.Size(251, 35);
            this.tableLayoutPanelGameControl.TabIndex = 4;
            // 
            // buttonRemoveGame
            // 
            this.buttonRemoveGame.Anchor = ((System.Windows.Forms.AnchorStyles)((((System.Windows.Forms.AnchorStyles.Top | System.Windows.Forms.AnchorStyles.Bottom) 
            | System.Windows.Forms.AnchorStyles.Left) 
            | System.Windows.Forms.AnchorStyles.Right)));
            this.buttonRemoveGame.Enabled = false;
            this.buttonRemoveGame.Location = new System.Drawing.Point(128, 3);
            this.buttonRemoveGame.Name = "buttonRemoveGame";
            this.buttonRemoveGame.Size = new System.Drawing.Size(120, 29);
            this.buttonRemoveGame.TabIndex = 4;
            this.buttonRemoveGame.Text = "Remove";
            this.buttonRemoveGame.UseVisualStyleBackColor = true;
            this.buttonRemoveGame.Click += new System.EventHandler(this.ButtonRemoveGame_Click);
            // 
            // buttonBackUp
            // 
            this.buttonBackUp.Anchor = ((System.Windows.Forms.AnchorStyles)((((System.Windows.Forms.AnchorStyles.Top | System.Windows.Forms.AnchorStyles.Bottom) 
            | System.Windows.Forms.AnchorStyles.Left) 
            | System.Windows.Forms.AnchorStyles.Right)));
            this.buttonBackUp.Location = new System.Drawing.Point(263, 439);
            this.buttonBackUp.Margin = new System.Windows.Forms.Padding(6);
            this.buttonBackUp.Name = "buttonBackUp";
            this.buttonBackUp.Size = new System.Drawing.Size(589, 29);
            this.buttonBackUp.TabIndex = 17;
            this.buttonBackUp.Text = "Back up";
            this.buttonBackUp.UseVisualStyleBackColor = true;
            this.buttonBackUp.Visible = false;
            this.buttonBackUp.Click += new System.EventHandler(this.ButtonBackUp_Click);
            // 
            // Form
            // 
            this.AutoScaleDimensions = new System.Drawing.SizeF(8F, 15F);
            this.AutoScaleMode = System.Windows.Forms.AutoScaleMode.Font;
            this.ClientSize = new System.Drawing.Size(882, 498);
            this.Controls.Add(this.tableLayoutPanelRoot);
            this.Name = "Form";
            this.ShowIcon = false;
            this.Text = "Game Save Manager";
            this.tableLayoutPanelRoot.ResumeLayout(false);
            this.tableLayoutPanelGameManager.ResumeLayout(false);
            this.tableLayoutPanelGameManager.PerformLayout();
            this.tableLayoutPanelGameControl.ResumeLayout(false);
            this.ResumeLayout(false);

        }

        private System.Windows.Forms.Label labelGames;
        private System.Windows.Forms.ListBox gameList;
        private System.Windows.Forms.Button buttonAddGame;
        private System.Windows.Forms.TableLayoutPanel tableLayoutPanelRoot;
        private System.Windows.Forms.TableLayoutPanel tableLayoutPanelGameManager;
        private System.Windows.Forms.Label labelBackupPath;
        private System.Windows.Forms.TextBox textBoxBackupPath;
        private System.Windows.Forms.Button buttonChooseBackupPath;
        private System.Windows.Forms.TableLayoutPanel tableLayoutPanelGameControl;
        private System.Windows.Forms.Button buttonRemoveGame;
        private System.Windows.Forms.Label labelGameName;
        private System.Windows.Forms.TextBox textBoxGameName;
        private System.Windows.Forms.Button buttonSaveName;
        private System.Windows.Forms.Button buttonOpenBackup;
        private System.Windows.Forms.Button buttonBackUp;
        private System.Windows.Forms.Label labelSaveRoot;
        private System.Windows.Forms.TextBox textBoxSaveRoot;
        private System.Windows.Forms.Button buttonChooseSaveRoot;
        private System.Windows.Forms.Button buttonOpenSaveRoot;
        private System.Windows.Forms.TreeView treeView;
    }
}

