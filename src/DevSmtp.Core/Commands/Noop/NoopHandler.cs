using DevSmtp.Core.Stores;

namespace DevSmtp.Core.Commands
{
    public sealed class NoopHandler : ICommandHandler<Noop, NoopResult>
    {
        private readonly IDataStore _dataStore;

        public NoopHandler(IDataStore dataStore)
        {
            this._dataStore = dataStore ?? throw new ArgumentNullException(nameof(dataStore));
        }

        public Task<NoopResult> ExecuteAsync(Noop command, CancellationToken cancellationToken = default)
        {
            throw new NotImplementedException();
        }
    }
}
