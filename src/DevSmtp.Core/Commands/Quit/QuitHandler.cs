using DevSmtp.Core.Stores;

namespace DevSmtp.Core.Commands
{
    public sealed class QuitHandler : ICommandHandler<Quit, QuitResult>
    {
        private readonly IDataStore _dataStore;

        public QuitHandler(IDataStore dataStore)
        {
            this._dataStore = dataStore ?? throw new ArgumentNullException(nameof(dataStore));
        }

        public Task<QuitResult> ExecuteAsync(Quit command, CancellationToken cancellationToken = default)
        {
            throw new NotImplementedException();
        }
    }
}
